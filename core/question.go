package core

import (
	"encoding/json"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	ytypes "github.com/yu-org/yu/core/types"
	"uask-chain/filestore"
	"uask-chain/search"
	"uask-chain/types"
)

type Question struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	sch       search.Search

	answer *Answer `tripod:"answer"`
}

func NewQuestion(fileStore filestore.FileStore, sch search.Search) *Question {
	tri := tripod.NewTripod()
	q := &Question{Tripod: tri, fileStore: fileStore, sch: sch}
	q.SetWritings(q.AddQuestion, q.UpdateQuestion, q.DeleteQuestion)
	q.SetTxnChecker(q)
	return q
}

func (q *Question) CheckTxn(txn *ytypes.SignedTxn) error {
	req := &types.QuestionAddRequest{}
	err := txn.BindJsonParams(req)
	if err != nil {
		return err
	}
	return checkOffchainOrStoreOnchain(txn.FromP2p(), req.Content, q.fileStore)
}

func (q *Question) AddQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	asker := ctx.GetCaller()
	req := &types.QuestionAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	scheme := &types.QuestionScheme{
		ID:          ctx.Txn.TxnHash.String(),
		Title:       req.Title,
		Asker:       asker,
		Tags:        req.Tags,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	err = q.setQuestionScheme(scheme)
	if err != nil {
		return err
	}
	// add search
	contentByt, err := q.fileStore.Get(req.Content.Hash)
	if err != nil {
		return err
	}
	err = q.sch.AddDoc(&types.Question{
		ID:          scheme.ID,
		Title:       scheme.Title,
		FileContent: contentByt,
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

	question, err := q.getQuestionScheme(req.ID)
	if err != nil {
		return err
	}
	if question.Asker != asker {
		return types.ErrNoPermission
	}

	scheme := &types.QuestionScheme{
		ID:          req.ID,
		Title:       req.Title,
		FileHash:    req.Content.Hash,
		Asker:       asker,
		Tags:        req.Tags,
		Timestamp:   req.Timestamp,
		Recommender: req.Recommender,
	}
	err = q.setQuestionScheme(scheme)
	if err != nil {
		return err
	}
	// update search
	contentByt, err := q.fileStore.Get(req.Content.Hash)
	if err != nil {
		return err
	}
	err = q.sch.UpdateDoc(scheme.ID, &types.Question{
		ID:          scheme.ID,
		Title:       scheme.Title,
		FileContent: contentByt,
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
	scheme, err := q.getQuestionScheme(id)
	if err != nil {
		return err
	}
	if asker != scheme.Asker {
		return types.ErrNoPermission
	}
	q.Delete([]byte(id))
	return q.sch.DeleteDoc(id)
}

func (q *Question) setQuestionScheme(scheme *types.QuestionScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	q.Set([]byte(scheme.ID), byt)
	return nil
}

func (q *Question) getQuestionScheme(id string) (*types.QuestionScheme, error) {
	byt, err := q.Get([]byte(id))
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
	return q.Exist([]byte(id))
}

func checkOffchainOrStoreOnchain(fromP2P bool, info *types.StoreInfo, store filestore.FileStore) error {
	if !fromP2P {
		// from RPC, store it into ipfs and clean the content.
		hash, err := store.Put("", info)
		if err != nil {
			return err
		}
		info.Hash = hash
		info.Url = store.Url()
		info.Content = nil
		return nil
	}
	// check ipfs if file exists
	byt, err := store.Get(info.Hash)
	if err != nil {
		return err
	}
	if byt == nil {
		return types.ErrFileNotfound
	}
	return nil
}
