package comment

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"gorm.io/gorm"
	"uask-chain/core/answer"
	"uask-chain/core/comment/orm"
	"uask-chain/core/question"
	"uask-chain/core/user"
	"uask-chain/filestore"
	"uask-chain/types"
)

type Comment struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	db        *orm.Database

	Question *question.Question `tripod:"question"`
	Answer   *answer.Answer     `tripod:"answer"`
	user     *user.User         `tripod:"user"`
}

func NewComment(fileStore filestore.FileStore, db *gorm.DB) *Comment {
	tri := tripod.NewTripod()
	database, err := orm.NewDB(db)
	if err != nil {
		logrus.Fatal("init comment db failed: ", err)
	}
	c := &Comment{Tripod: tri, fileStore: fileStore, db: database}
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

	err = c.user.CheckReputation(commenter, types.AddCommentReputationNeed)
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

	scheme := &orm.CommentScheme{
		ID:        ctx.GetTxnHash().String(),
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

	scheme := &orm.CommentScheme{
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

func (c *Comment) setCommentState(scheme *orm.CommentScheme) error {
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
		if !c.Answer.ExistAnswer(answerID) {
			return types.ErrAnswerNotFound
		}
	}
	if questionID != "" {
		if !c.Question.ExistQuestion(questionID) {
			return types.ErrQuestionNotFound
		}
	}
	return nil
}
