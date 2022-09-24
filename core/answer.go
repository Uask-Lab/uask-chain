package core

import (
	"encoding/json"
	"fmt"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/core/types"
	"uask-chain/filestore"
)

type Answer struct {
	*tripod.Tripod
	fileStore filestore.FileStore
}

func NewAnswer(fileStore filestore.FileStore) *Answer {
	tri := tripod.NewTripod("answer")
	a := &Answer{Tripod: tri, fileStore: fileStore}
	a.SetExec(a.AddAnswer).SetExec(a.UpdateAnswer)
	return a
}

func (a *Answer) AddAnswer(ctx *context.Context) error {
	ctx.SetLei(50)

	answerer := ctx.Caller
	req := &types.AnswerAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	// check if question exists
	q := a.GetTripod("question").(*Question)
	if !q.existQuestion(req.QID) {
		return ErrQuestionNotFound
	}

	id := fmt.Sprintf("%s%s%s", answerer.String(), req.QID, req.Timestamp)
	stub, err := a.fileStore.Put(id, req.Content)
	if err != nil {
		return err
	}

	scheme := &types.AnswerScheme{
		ID:          id,
		QID:         req.QID,
		Answerer:    answerer,
		ContentStub: stub,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	a.State.Set(a, []byte(id), byt)
	ctx.EmitEvent(fmt.Sprintf("add answer(%s) to question(%s) successfully by answerer(%s)!", scheme.ID, scheme.QID, answerer.String()))
	return nil
}

func (a *Answer) UpdateAnswer(ctx *context.Context) error {
	ctx.SetLei(50)

	answerer := ctx.Caller
	req := &types.AnswerUpdateRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	if !a.existAnswer(req.ID) {
		return ErrAnswerNotFound
	}

	answer, err := a.getAnswer(req.ID)
	if err != nil {
		return err
	}
	if answer.Answerer != answerer {
		return ErrNoPermission
	}

	stub, err := a.fileStore.Put(req.ID, req.Content)
	if err != nil {
		return err
	}
	scheme := &types.AnswerScheme{
		ID:          req.ID,
		QID:         req.QID,
		Answerer:    answerer,
		ContentStub: stub,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	a.State.Set(a, []byte(req.ID), byt)
	ctx.EmitEvent(fmt.Sprintf("update answer(%s) successfully!", req.ID))
	return nil

}

func (a *Answer) getAnswer(id string) (*types.AnswerScheme, error) {
	byt, err := a.State.Get(a, []byte(id))
	if err != nil {
		return nil, err
	}
	scheme := &types.AnswerScheme{}
	err = json.Unmarshal(byt, scheme)
	if err != nil {
		return nil, err
	}
	return scheme, nil
}

func (a *Answer) existAnswer(id string) bool {
	return a.State.Exist(a, []byte(id))
}
