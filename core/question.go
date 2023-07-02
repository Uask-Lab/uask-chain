package core

import (
	"encoding/json"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/db"
	"uask-chain/filestore"
	"uask-chain/search"
	"uask-chain/types"
)

type Question struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	sch       search.Search
	db        *db.Database

	answer *Answer `tripod:"answer"`
}

func NewQuestion(fileStore filestore.FileStore, sch search.Search, db *db.Database) *Question {
	tri := tripod.NewTripod()
	q := &Question{Tripod: tri, fileStore: fileStore, sch: sch, db: db}
	q.SetWritings(q.AddQuestion, q.UpdateQuestion, q.DeleteQuestion)
	q.SetReadings(q.GetQuestion, q.SearchQuestion)
	return q
}

func (q *Question) GetQuestion(ctx *context.ReadContext) error {
	sch, err := q.db.GetQuestion(ctx.GetString("id"))
	if err != nil {
		return err
	}
	fileByt, err := q.fileStore.Get(sch.FileHash)
	if err != nil {
		return err
	}
	question := &types.QuestionInfo{
		QuestionDoc: types.QuestionDoc{
			ID:          sch.ID,
			Title:       sch.Title,
			Content:     fileByt,
			Asker:       sch.Asker,
			Tags:        sch.Tags,
			Timestamp:   sch.Timestamp,
			Recommender: sch.Recommender,
		},
	}
	return ctx.Json(question)
}

func (q *Question) SearchQuestion(ctx *context.ReadContext) error {
	phrase := ctx.GetString("phrase")
	results, err := q.sch.SearchDoc(phrase)
	if err != nil {
		return err
	}
	return ctx.Json(results)
}

func (q *Question) AddQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	asker := ctx.GetCaller()
	req := &types.QuestionAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	fileHash, err := q.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.QuestionScheme{
		ID:          ctx.Txn.TxnHash.String(),
		Title:       req.Title,
		Asker:       asker,
		FileHash:    fileHash,
		Tags:        req.Tags,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	err = q.setQuestionState(scheme)
	if err != nil {
		return err
	}

	// store into database
	err = q.db.AddQuestion(scheme)
	if err != nil {
		return err
	}

	// add search
	err = q.sch.AddDoc(&types.QuestionDoc{
		ID:          scheme.ID,
		Title:       scheme.Title,
		Content:     req.Content,
		Asker:       scheme.Asker,
		Tags:        scheme.Tags,
		Timestamp:   scheme.Timestamp,
		Recommender: scheme.Recommender,
	})
	if err != nil {
		return err
	}

	ctx.EmitStringEvent("add question(%s) successfully by asker(%s)! question-id=%s", scheme.Title, asker.String(), scheme.ID)
	return nil
}

func (q *Question) UpdateQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	asker := ctx.GetCaller()
	req := &types.QuestionUpdateRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	question, err := q.db.GetQuestion(req.ID)
	if err != nil {
		return err
	}
	if question.Asker != asker {
		return types.ErrNoPermission
	}

	// remove old answer and store new one.
	err = q.fileStore.Remove(question.FileHash)
	if err != nil {
		return err
	}
	fileHash, err := q.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.QuestionScheme{
		ID:          req.ID,
		Title:       req.Title,
		FileHash:    fileHash,
		Asker:       asker,
		Tags:        req.Tags,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	err = q.setQuestionState(scheme)
	if err != nil {
		return err
	}

	// update database
	err = q.db.UpdateQuestion(scheme)
	if err != nil {
		return err
	}

	// update doc
	err = q.sch.UpdateDoc(scheme.ID, &types.QuestionDoc{
		ID:          scheme.ID,
		Title:       scheme.Title,
		Content:     req.Content,
		Asker:       scheme.Asker,
		Tags:        scheme.Tags,
		Timestamp:   scheme.Timestamp,
		Recommender: scheme.Recommender,
	})
	if err != nil {
		return err
	}

	ctx.EmitStringEvent("update question(%s) successfully!", req.ID)
	return nil
}

func (q *Question) DeleteQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)
	id := ctx.GetString("id")
	asker := ctx.GetCaller()
	scheme, err := q.db.GetQuestion(id)
	if err != nil {
		return err
	}
	if asker != scheme.Asker {
		return types.ErrNoPermission
	}
	q.Delete([]byte(id))
	err = q.db.DeleteQuestion(id)
	if err != nil {
		return err
	}
	return q.sch.DeleteDoc(id)
}

func (q *Question) setQuestionState(scheme *types.QuestionScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	hashByt := common.Sha256(byt)

	q.Set([]byte(scheme.ID), hashByt)
	return nil
}

func (q *Question) existQuestion(id string) bool {
	return q.Exist([]byte(id))
}
