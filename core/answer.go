package core

import (
	"encoding/json"
	"fmt"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	ytypes "github.com/yu-org/yu/core/types"
	"uask-chain/filestore"
	"uask-chain/types"
)

type Answer struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	Question  *Question `tripod:"question"`
}

func NewAnswer(fileStore filestore.FileStore) *Answer {
	tri := tripod.NewTripod("answer")
	a := &Answer{Tripod: tri, fileStore: fileStore}
	a.SetWritings(a.AddAnswer, a.UpdateAnswer)
	a.SetTxnChecker(a)
	return a
}

func (a *Answer) CheckTxn(txn *ytypes.SignedTxn) error {
	req := &types.AnswerAddRequest{}
	err := txn.BindJsonParams(req)
	if err != nil {
		return err
	}
	return checkOffchainStore(req.Content, a.fileStore)
}

func (a *Answer) AddAnswer(ctx *context.WriteContext) error {
	ctx.SetLei(50)

	answerer := ctx.GetCaller()
	req := &types.AnswerAddRequest{}
	err := ctx.Bindjson(req)
	if err != nil {
		return err
	}

	// check if question exists
	if !a.Question.existQuestion(req.QID) {
		return types.ErrQuestionNotFound
	}

	scheme := &types.AnswerScheme{
		ID:          ctx.Txn.TxnHash.String(),
		QID:         req.QID,
		Answerer:    answerer,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	err = a.setAnswer(scheme)
	if err != nil {
		return err
	}
	return ctx.EmitEvent(fmt.Sprintf("add answer(%s) to question(%s) successfully by answerer(%s)!", scheme.ID, scheme.QID, answerer.String()))
}

func (a *Answer) UpdateAnswer(ctx *context.WriteContext) error {
	ctx.SetLei(50)

	answerer := ctx.GetCaller()
	req := &types.AnswerUpdateRequest{}
	err := ctx.Bindjson(req)
	if err != nil {
		return err
	}

	if !a.existAnswer(req.ID) {
		return types.ErrAnswerNotFound
	}

	answer, err := a.getAnswer(req.ID)
	if err != nil {
		return err
	}
	if answer.Answerer != answerer {
		return types.ErrNoPermission
	}

	scheme := &types.AnswerScheme{
		ID:          req.ID,
		QID:         req.QID,
		Answerer:    answerer,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	err = a.setAnswer(scheme)
	if err != nil {
		return err
	}
	return ctx.EmitEvent(fmt.Sprintf("update answer(%s) successfully!", req.ID))
}

func (a *Answer) setAnswer(scheme *types.AnswerScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	a.State.Set(a, []byte(scheme.ID), byt)
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
