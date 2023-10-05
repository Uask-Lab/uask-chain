package answer

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"gorm.io/gorm"
	"uask-chain/core/answer/orm"
	"uask-chain/core/question"
	"uask-chain/core/user"
	"uask-chain/filestore"
	"uask-chain/types"
)

type Answer struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	db        *orm.Database

	Question *question.Question `tripod:"question"`
	user     *user.User         `tripod:"user"`
}

func NewAnswer(fileStore filestore.FileStore, db *gorm.DB) *Answer {
	tri := tripod.NewTripod()
	database, err := orm.NewDB(db)
	if err != nil {
		logrus.Fatal("init answer db failed: ", err)
	}
	a := &Answer{Tripod: tri, fileStore: fileStore, db: database}
	a.SetWritings(
		a.AddAnswer,
		a.UpdateAnswer,
		a.DeleteAnswer,
		a.UpVote,
		a.DownVote,
		a.PickUp,
		a.Drop,
	)
	a.SetReadings(a.GetAnswer)
	return a
}

func (a *Answer) AddAnswer(ctx *context.WriteContext) error {
	ctx.SetLei(50)

	answerer := ctx.GetCaller()
	req := &types.AnswerAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	err = a.user.CheckReputation(answerer, types.AddAnswerReputationNeed)
	if err != nil {
		return err
	}

	// check if question exists
	if !a.Question.ExistQuestion(req.QID) {
		return types.ErrQuestionNotFound
	}

	fileHash, err := a.fileStore.Put([]byte(req.Content))
	if err != nil {
		return err
	}

	scheme := &orm.AnswerScheme{
		ID:        ctx.GetTxnHash().String(),
		QID:       req.QID,
		FileHash:  fileHash,
		Answerer:  answerer.String(),
		Timestamp: int64(ctx.GetTimestamp()),
	}
	err = a.setAnswerState(scheme)
	if err != nil {
		return err
	}

	// store into database
	err = a.db.AddAnswer(scheme)
	if err != nil {
		return err
	}

	return ctx.EmitJsonEvent(map[string]string{
		"writing":     "add_answer",
		"id":          scheme.ID,
		"question_id": scheme.QID,
		"answerer":    answerer.String(),
	})
}

func (a *Answer) UpdateAnswer(ctx *context.WriteContext) error {
	ctx.SetLei(50)

	answerer := ctx.GetCaller()
	req := &types.AnswerUpdateRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	if !a.ExistAnswer(req.ID) {
		return types.ErrAnswerNotFound
	}

	answer, err := a.db.GetAnswer(req.ID)
	if err != nil {
		return err
	}
	if answer.Answerer != answerer.String() {
		return types.ErrNoPermission
	}

	// remove old answer and store new one.
	err = a.fileStore.Remove(answer.FileHash)
	if err != nil {
		return err
	}
	fileHash, err := a.fileStore.Put([]byte(req.Content))
	if err != nil {
		return err
	}

	scheme := &orm.AnswerScheme{
		ID:        req.ID,
		QID:       req.QID,
		FileHash:  fileHash,
		Answerer:  answerer.String(),
		Timestamp: int64(ctx.GetTimestamp()),
	}
	err = a.setAnswerState(scheme)
	if err != nil {
		return err
	}

	// update database
	err = a.db.UpdateAnswer(scheme)
	if err != nil {
		return err
	}

	return ctx.EmitJsonEvent(map[string]string{"writing": "update_answer", "id": req.ID})
}

func (a *Answer) GetAnswer(ctx *context.ReadContext) {
	id := ctx.GetString("id")
	scheme, err := a.db.GetAnswer(id)
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}
	fileByt, err := a.fileStore.Get(scheme.FileHash)
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}

	answer := &types.AnswerInfo{
		ID:        scheme.ID,
		QID:       scheme.QID,
		Content:   string(fileByt),
		Answerer:  scheme.Answerer,
		Timestamp: scheme.Timestamp,
	}
	ctx.JsonOk(types.Ok(answer))
}

func (a *Answer) DeleteAnswer(ctx *context.WriteContext) error {
	ctx.SetLei(10)
	id := ctx.GetString("id")
	answerer := ctx.GetCaller()
	scheme, err := a.db.GetAnswer(id)
	if err == types.ErrAnswerNotFound {
		return ctx.EmitJsonEvent(map[string]string{"writing": "delete_answer", "id": id, "status": "none"})
	}
	if err != nil {
		return err
	}
	if answerer.String() != scheme.Answerer {
		return types.ErrNoPermission
	}
	a.Delete([]byte(id))
	err = a.db.DeleteAnswer(id)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"writing": "delete_answer", "id": id, "status": "success"})
}

func (a *Answer) setAnswerState(scheme *orm.AnswerScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	a.Set([]byte(scheme.ID), byt)
	return nil
}

func (a *Answer) ExistAnswer(id string) bool {
	return a.Exist([]byte(id))
}
