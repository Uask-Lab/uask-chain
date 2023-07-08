package core

import (
	"encoding/json"
	"github.com/yu-org/yu/common"
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

	Answer *Answer `tripod:"answer"`
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

	err = c.ifReplyExist(req.AID, req.CID)
	if err != nil {
		return err
	}

	fileHash, err := c.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:        ctx.Txn.TxnHash.String(),
		QID:       req.QID,
		AID:       req.AID,
		CID:       req.CID,
		FileHash:  fileHash,
		Commenter: commenter.String(),
		Timestamp: req.Timestamp,
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

	err = c.ifReplyExist(req.AID, req.CID)
	if err != nil {
		return err
	}

	// remove old answer and store new one.
	err = c.fileStore.Remove(comment.FileHash)
	if err != nil {
		return err
	}
	fileHash, err := c.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:        req.ID,
		AID:       req.AID,
		CID:       req.CID,
		FileHash:  fileHash,
		Commenter: commenter.String(),
		Timestamp: req.Timestamp,
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
		return nil
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
	return ctx.EmitJsonEvent(map[string]string{"writing": "delete_comment", "id": id})
}

func (c *Comment) GetComment(ctx *context.ReadContext) error {
	sch, err := c.db.GetComment(ctx.GetString("id"))
	if err != nil {
		return err
	}
	fileByt, err := c.fileStore.Get(sch.FileHash)
	if err != nil {
		return err
	}
	comment := &types.CommentInfo{
		CommentUpdateRequest: types.CommentUpdateRequest{
			ID: sch.ID,
			CommentAddRequest: types.CommentAddRequest{
				QID:       sch.QID,
				AID:       sch.AID,
				CID:       sch.CID,
				Content:   fileByt,
				Timestamp: sch.Timestamp,
			},
		},
	}
	return ctx.Json(comment)
}

func (c *Comment) setCommentState(scheme *types.CommentScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	hashByt := common.Sha256(byt)

	c.Set([]byte(scheme.ID), hashByt)
	return nil
}

func (c *Comment) existComment(id string) bool {
	return c.Exist([]byte(id))
}

func (c *Comment) ifReplyExist(answerID, commentID string) error {
	if answerID != "" {
		if !c.Answer.existAnswer(answerID) {
			return types.ErrAnswerNotFound
		}
	}
	if commentID != "" {
		if !c.existComment(commentID) {
			return types.ErrCommentNotFound
		}
	}
	return nil
}
