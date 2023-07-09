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
	q.SetReadings(q.ListQuestions, q.GetQuestion, q.SearchQuestion)
	return q
}

func (q *Question) ListQuestions(ctx *context.ReadContext) error {
	pageSize := ctx.GetInt("pageSize")
	page := ctx.GetInt("page")

	qschs, err := q.db.ListQuestions(pageSize, (page-1)*pageSize)
	if err != nil {
		return err
	}

	var infos []*types.QuestionInfo
	for _, qsch := range qschs {
		info, serr := q.scheme2Info(qsch)
		if serr != nil {
			return serr
		}
		infos = append(infos, info)
	}

	return ctx.Json(infos)
}

func (q *Question) GetQuestion(ctx *context.ReadContext) error {
	sch, err := q.db.GetQuestion(ctx.GetString("id"))
	if err != nil {
		return err
	}
	question, err := q.scheme2Info(sch)
	if err != nil {
		return err
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
		Asker:       asker.String(),
		FileHash:    fileHash,
		Tags:        req.Tags,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender.String(),
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
		Asker:       common.HexToAddress(scheme.Asker),
		Tags:        scheme.Tags,
		Timestamp:   scheme.Timestamp,
		Recommender: common.HexToAddress(scheme.Recommender),
	})
	if err != nil {
		return err
	}

	return ctx.EmitJsonEvent(map[string]string{
		"writing": "add_question",
		"id":      scheme.ID,
		"title":   scheme.Title,
		"asker":   asker.String(),
	})
}

func (q *Question) UpdateQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	asker := ctx.GetCaller()
	req := &types.QuestionUpdateRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	if !q.existQuestion(req.ID) {
		return types.ErrQuestionNotFound
	}

	question, err := q.db.GetQuestion(req.ID)
	if err != nil {
		return err
	}
	if question.Asker != asker.String() {
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
		Asker:       asker.String(),
		Tags:        req.Tags,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender.String(),
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
		Asker:       common.HexToAddress(scheme.Asker),
		Tags:        scheme.Tags,
		Timestamp:   scheme.Timestamp,
		Recommender: common.HexToAddress(scheme.Recommender),
	})
	if err != nil {
		return err
	}

	return ctx.EmitJsonEvent(map[string]string{"writing": "update_question", "id": scheme.ID})
}

func (q *Question) DeleteQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)
	id := ctx.GetString("id")
	asker := ctx.GetCaller()
	scheme, err := q.db.GetQuestion(id)
	if err == types.ErrQuestionNotFound {
		return ctx.EmitJsonEvent(map[string]string{"writing": "delete_question", "id": id, "status": "none"})
	}
	if err != nil {
		return err
	}
	if asker.String() != scheme.Asker {
		return types.ErrNoPermission
	}
	q.Delete([]byte(id))
	err = q.db.DeleteQuestion(id)
	if err != nil {
		return err
	}
	err = q.sch.DeleteDoc(id)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"writing": "delete_question", "id": id, "status": "success"})
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

func (q *Question) scheme2Info(sch *types.QuestionScheme) (*types.QuestionInfo, error) {
	fileByt, err := q.fileStore.Get(sch.FileHash)
	if err != nil {
		return nil, err
	}
	return &types.QuestionInfo{
		QuestionDoc: types.QuestionDoc{
			ID:          sch.ID,
			Title:       sch.Title,
			Content:     fileByt,
			Asker:       common.HexToAddress(sch.Asker),
			Tags:        sch.Tags,
			Timestamp:   sch.Timestamp,
			Recommender: common.HexToAddress(sch.Recommender),
		},
	}, nil
}
