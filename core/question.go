package core

import (
	"encoding/json"
	"fmt"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/filestore"
)

type Question struct {
	*tripod.Tripod
	fileStore filestore.FileStore
}

func NewQuestion(fileStore filestore.FileStore) *Question {
	tri := tripod.NewTripod("question")
	q := &Question{Tripod: tri, fileStore: fileStore}
	q.SetExec(q.AddQuestion).SetExec(q.UpdateQuestion)
	q.SetQueries(q.QueryQuestion)
	return q
}

func (q *Question) AddQuestion(ctx *context.Context) error {
	ctx.SetLei(10)

	asker := ctx.Caller
	req := &QuestionRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	// TODO: Lock the amount of balance for reward.

	id := fmt.Sprintf("%s%s%s", asker.String(), req.Title, req.Timestamp)
	stub, err := q.fileStore.Put(id, req.Content)
	if err != nil {
		return err
	}

	scheme := &QuestionScheme{
		ID:          id,
		Title:       req.Title,
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
	q.State.Set(q, asker.Bytes(), byt)
	ctx.EmitEvent(fmt.Sprintf("add question(%s) successfully by asker(%s)! question-id=%s", scheme.Title, asker.String(), scheme.ID))
	return nil
}

func (q *Question) UpdateQuestion(ctx *context.Context) error {

}

func (q *Question) QueryQuestion(ctx *context.Context) (interface{}, error) {

}
