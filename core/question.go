package core

import (
	"encoding/json"
	"fmt"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/filestore"
	"uask-chain/types"
)

type Question struct {
	*tripod.Tripod
	fileStore filestore.FileStore
}

func NewQuestion(fileStore filestore.FileStore) *Question {
	tri := tripod.NewTripod("question")
	q := &Question{Tripod: tri, fileStore: fileStore}
	q.SetExec(q.AddQuestion).SetExec(q.UpdateQuestion)
	return q
}

func (q *Question) AddQuestion(ctx *context.Context) error {
	ctx.SetLei(10)

	asker := ctx.Caller
	req := ctx.ParamsValue.(*types.QuestionAddRequest)

	// TODO: Lock the amount of balance for reward.

	id := fmt.Sprintf("%s%s%s", asker.String(), req.Title, req.Timestamp)
	stub, err := q.fileStore.Put(id, req.Content)
	if err != nil {
		return err
	}

	scheme := &types.QuestionScheme{
		ID:          id,
		Title:       req.Title,
		Asker:       asker,
		ContentStub: stub,
		Tags:        req.Tags,
		Reward:      req.Reward,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	q.State.Set(q, []byte(id), byt)
	ctx.EmitEvent(fmt.Sprintf("add question(%s) successfully by asker(%s)! question-id=%s", scheme.Title, asker.String(), scheme.ID))
	return nil
}

func (q *Question) UpdateQuestion(ctx *context.Context) error {
	ctx.SetLei(10)

	asker := ctx.Caller
	req := ctx.ParamsValue.(*types.QuestionUpdateRequest)

	question, err := q.getQuestion(req.ID)
	if err != nil {
		return err
	}
	if question.Asker != asker {
		return ErrNoPermission
	}

	// TODO: Lock the amount of balance for reward.

	stub, err := q.fileStore.Put(req.ID, req.Content)
	if err != nil {
		return err
	}

	scheme := &types.QuestionScheme{
		ID:          req.ID,
		Title:       req.Title,
		Asker:       asker,
		ContentStub: stub,
		Tags:        req.Tags,
		Reward:      req.Reward,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	q.State.Set(q, []byte(req.ID), byt)
	ctx.EmitEvent(fmt.Sprintf("update question(%s) successfully!", req.ID))
	return nil
}

func (q *Question) getQuestion(id string) (*types.QuestionScheme, error) {
	byt, err := q.State.Get(q, []byte(id))
	if err != nil {
		return nil, err
	}
	scheme := &types.QuestionScheme{}
	err = json.Unmarshal(byt, scheme)
	if err != nil {
		return nil, err
	}
	return scheme, nil
}

func (q *Question) existQuestion(id string) bool {
	return q.State.Exist(q, []byte(id))
}
