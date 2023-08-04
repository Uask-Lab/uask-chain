package core

import (
	"encoding/json"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/db"
	"uask-chain/filestore"
	"uask-chain/types"
)

type Comment struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	db        *db.Database

	Question *Question `tripod:"question"`
	Answer   *Answer   `tripod:"answer"`
}

func NewComment(fileStore filestore.FileStore, db *db.Database) *Comment {
	tri := tripod.NewTripod()
	c := &Comment{Tripod: tri, fileStore: fileStore, db: db}
	c.SetWritings(c.AddComment, c.UpdateComment, c.DeleteComment)
	c.SetReadings(c.GetComment)
	return c
}

func (c *Comment) AddComment(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	commenter := ctx.GetCaller()
	req := &types.CommentAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	err = c.ifReplyExist(req.QID, req.AID)
	if err != nil {
		return err
	}

	fileHash, err := c.fileStore.Put([]byte(req.Content))
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:        ctx.Txn.TxnHash.String(),
		QID:       req.QID,
		AID:       req.AID,
		FileHash:  fileHash,
		Commenter: commenter.String(),
		Timestamp: int64(ctx.GetTimestamp()),
	}
	err = c.setCommentState(scheme)
	if err != nil {
		return err
	}

	// store into database
	err = c.db.AddComment(scheme)
	if err != nil {
		return err
	}

	return ctx.EmitJsonEvent(map[string]string{
		"writing":   "add_comment",
		"id":        scheme.ID,
		"commenter": commenter.String(),
	})
}

func (c *Comment) UpdateComment(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	commenter := ctx.GetCaller()
	req := &types.CommentUpdateRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	comment, err := c.db.GetComment(req.ID)
	if err != nil {
		return err
	}
	if comment.Commenter != commenter.String() {
		return types.ErrNoPermission
	}

	err = c.ifReplyExist(req.QID, req.AID)
	if err != nil {
		return err
	}

	// remove old answer and store new one.
	err = c.fileStore.Remove(comment.FileHash)
	if err != nil {
		return err
	}
	fileHash, err := c.fileStore.Put([]byte(req.Content))
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:        req.ID,
		AID:       req.AID,
		FileHash:  fileHash,
		Commenter: commenter.String(),
		Timestamp: int64(ctx.GetTimestamp()),
	}
	err = c.setCommentState(scheme)
	if err != nil {
		return err
	}

	// update database
	err = c.db.UpdateComment(scheme)
	if err != nil {
		return err
	}

	return ctx.EmitJsonEvent(map[string]string{"writing": "update_comment", "id": req.ID})
}

func (c *Comment) DeleteComment(ctx *context.WriteContext) error {
	ctx.SetLei(10)
	id := ctx.GetString("id")
	commenter := ctx.GetCaller()
	scheme, err := c.db.GetComment(id)
	if err == types.ErrCommentNotFound {
		return ctx.EmitJsonEvent(map[string]string{"writing": "delete_comment", "id": id, "status": "none"})
	}
	if err != nil {
		return err
	}
	if commenter.String() != scheme.Commenter {
		return types.ErrNoPermission
	}
	c.Delete([]byte(id))
	err = c.db.DeleteComment(id)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"writing": "delete_comment", "id": id, "status": "success"})
}

func (c *Comment) GetComment(ctx *context.ReadContext) {
	sch, err := c.db.GetComment(ctx.GetString("id"))
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}
	fileByt, err := c.fileStore.Get(sch.FileHash)
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}
	comment := &types.CommentInfo{
		ID:        sch.ID,
		QID:       sch.QID,
		AID:       sch.AID,
		Content:   string(fileByt),
		Commenter: sch.Commenter,
		Timestamp: sch.Timestamp,
	}
	ctx.JsonOk(types.Ok(comment))
}

func (c *Comment) setCommentState(scheme *types.CommentScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	c.Set([]byte(scheme.ID), byt)
	return nil
}

func (c *Comment) existComment(id string) bool {
	return c.Exist([]byte(id))
}

func (c *Comment) ifReplyExist(questionID, answerID string) error {
	if questionID == "" && answerID == "" {
		return types.ErrNoneToReply
	}
	if answerID != "" {
		if !c.Answer.existAnswer(answerID) {
			return types.ErrAnswerNotFound
		}
	}
	if questionID != "" {
		if !c.Question.existQuestion(questionID) {
			return types.ErrQuestionNotFound
		}
	}
	return nil
}
