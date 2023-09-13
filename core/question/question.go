package question

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"gorm.io/gorm"
	"uask-chain/core/question/orm"
	"uask-chain/filestore"
	"uask-chain/search"
	"uask-chain/types"
)

type Question struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	sch       search.Search
	db        *orm.Database
}

func NewQuestion(fileStore filestore.FileStore, sch search.Search, db *gorm.DB) *Question {
	database, err := orm.NewDB(db)
	if err != nil {
		logrus.Fatal("init question db failed: ", err)
	}
	tri := tripod.NewTripod()
	q := &Question{Tripod: tri, fileStore: fileStore, sch: sch, db: database}
	q.SetWritings(q.AddQuestion, q.UpdateQuestion, q.DeleteQuestion)
	q.SetReadings(q.ListQuestions, q.GetQuestion, q.SearchQuestion)
	return q
}

func (q *Question) ListQuestions(ctx *context.ReadContext) {
	pageSize := ctx.GetInt("pageSize")
	page := ctx.GetInt("page")

	qschs, err := q.db.ListQuestions(pageSize, (page-1)*pageSize)
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}

	var infos []*types.QuestionInfo
	for _, qsch := range qschs {
		info, serr := q.scheme2Info(qsch)
		if serr != nil {
			ctx.JsonOk(types.Error(serr))
			return
		}
		infos = append(infos, info)
	}

	ctx.JsonOk(types.Ok(infos))
}

func (q *Question) GetQuestion(ctx *context.ReadContext) {
	sch, err := q.db.GetQuestion(ctx.GetString("id"))
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}
	question, err := q.scheme2Info(sch)
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}
	ctx.JsonOk(types.Ok(question))
}

func (q *Question) SearchQuestion(ctx *context.ReadContext) {
	phrase := ctx.GetString("phrase")
	results, err := q.sch.SearchDoc(phrase)
	if err != nil {
		ctx.JsonOk(types.Error(err))
		return
	}
	ctx.JsonOk(types.Ok(results))
}

func (q *Question) AddQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	asker := ctx.GetCaller()
	req := &types.QuestionAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	fileHash, err := q.fileStore.Put([]byte(req.Content))
	if err != nil {
		return err
	}

	scheme := &orm.QuestionScheme{
		ID:        ctx.GetTxnHash().String(),
		Title:     req.Title,
		Asker:     asker.String(),
		FileHash:  fileHash,
		Tags:      req.Tags,
		Timestamp: int64(ctx.GetTimestamp()),
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
		ID:        scheme.ID,
		Title:     scheme.Title,
		Content:   req.Content,
		Asker:     scheme.Asker,
		Tags:      scheme.Tags,
		Timestamp: scheme.Timestamp,
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

	if !q.ExistQuestion(req.ID) {
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
	fileHash, err := q.fileStore.Put([]byte(req.Content))
	if err != nil {
		return err
	}

	scheme := &orm.QuestionScheme{
		ID:        req.ID,
		Title:     req.Title,
		FileHash:  fileHash,
		Asker:     asker.String(),
		Tags:      req.Tags,
		Timestamp: int64(ctx.GetTimestamp()),
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
		ID:        scheme.ID,
		Title:     scheme.Title,
		Content:   req.Content,
		Asker:     scheme.Asker,
		Tags:      scheme.Tags,
		Timestamp: scheme.Timestamp,
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

func (q *Question) setQuestionState(scheme *orm.QuestionScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	q.Set([]byte(scheme.ID), byt)
	return nil
}

func (q *Question) ExistQuestion(id string) bool {
	return q.Exist([]byte(id))
}

func (q *Question) scheme2Info(sch *orm.QuestionScheme) (*types.QuestionInfo, error) {
	fileByt, err := q.fileStore.Get(sch.FileHash)
	if err != nil {
		return nil, err
	}
	return &types.QuestionInfo{
		QuestionDoc: types.QuestionDoc{
			ID:        sch.ID,
			Title:     sch.Title,
			Content:   string(fileByt),
			Asker:     sch.Asker,
			Tags:      sch.Tags,
			Timestamp: sch.Timestamp,
		},
	}, nil
}
