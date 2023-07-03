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

type Answer struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	db        *db.Database
	// sch       search.Search

	Question *Question `tripod:"question"`
}

func NewAnswer(fileStore filestore.FileStore, db *db.Database) *Answer {
	tri := tripod.NewTripod()
	a := &Answer{Tripod: tri, fileStore: fileStore, db: db}
	a.SetWritings(a.AddAnswer, a.UpdateAnswer, a.DeleteAnswer)
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

	// check if question exists
	if !a.Question.existQuestion(req.QID) {
		return types.ErrQuestionNotFound
	}

	fileHash, err := a.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.AnswerScheme{
		ID:          ctx.Txn.TxnHash.String(),
		QID:         req.QID,
		FileHash:    fileHash,
		Answerer:    answerer.String(),
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender.String(),
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

	// add content into search
	//contentByt, err := a.fileStore.Get(req.Content.Hash)
	//if err != nil {
	//	return err
	//}
	//err = a.sch.AddDoc(&types.Answer{
	//	ID:          scheme.ID,
	//	Answerer:    scheme.Answerer,
	//	Content: contentByt,
	//	Recommender: scheme.Recommender,
	//	Timestamp:   scheme.Timestamp,
	//})
	//if err != nil {
	//	return err
	//}

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

	if !a.existAnswer(req.ID) {
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
	fileHash, err := a.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.AnswerScheme{
		ID:          req.ID,
		QID:         req.QID,
		FileHash:    fileHash,
		Answerer:    answerer.String(),
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender.String(),
	}
	err = a.setAnswerState(scheme)
	if err != nil {
		return err
	}

	// update database
	err = a.db.UpdateAnswer(*scheme)
	if err != nil {
		return err
	}

	// update content into search
	//contentByt, err := a.fileStore.Get(req.Content.Hash)
	//if err != nil {
	//	return err
	//}
	//err = a.sch.UpdateDoc(req.ID, &types.Answer{
	//	ID:          scheme.ID,
	//	Answerer:    scheme.Answerer,
	//	Content: contentByt,
	//	Recommender: scheme.Recommender,
	//	Timestamp:   scheme.Timestamp,
	//})
	//if err != nil {
	//	return err
	//}

	return ctx.EmitJsonEvent(map[string]string{"writing": "update_answer", "id": req.ID})
}

func (a *Answer) GetAnswer(ctx *context.ReadContext) error {
	id := ctx.GetString("id")
	scheme, err := a.db.GetAnswer(id)
	if err != nil {
		return err
	}
	fileByt, err := a.fileStore.Get(scheme.FileHash)
	if err != nil {
		return err
	}
	answer := &types.AnswerInfo{
		AnswerUpdateRequest: types.AnswerUpdateRequest{
			ID: scheme.ID,
			AnswerAddRequest: types.AnswerAddRequest{
				QID:         scheme.QID,
				Content:     fileByt,
				Timestamp:   scheme.Timestamp,
				Recommender: common.HexToAddress(scheme.Recommender),
			},
		},
	}
	return ctx.Json(answer)
}

func (a *Answer) DeleteAnswer(ctx *context.WriteContext) error {
	ctx.SetLei(10)
	id := ctx.GetString("id")
	answerer := ctx.GetCaller()
	scheme, err := a.db.GetAnswer(id)
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
	return ctx.EmitJsonEvent(map[string]string{"writing": "delete_answer", "id": id})
}

func (a *Answer) setAnswerState(scheme *types.AnswerScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	hashByt := common.Sha256(byt)

	a.Set([]byte(scheme.ID), hashByt)
	return nil
}

func (a *Answer) existAnswer(id string) bool {
	return a.Exist([]byte(id))
}
