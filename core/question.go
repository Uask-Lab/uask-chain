package core

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/asset"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	ytypes "github.com/yu-org/yu/core/types"
	"math/big"
	"uask-chain/filestore"
	"uask-chain/types"
)

type Question struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	asset     *asset.Asset `tripod:"asset"`
	answer    *Answer      `tripod:"answer"`
}

func NewQuestion(fileStore filestore.FileStore) *Question {
	tri := tripod.NewTripod("question")
	q := &Question{Tripod: tri, fileStore: fileStore}
	q.SetWritings(q.AddQuestion, q.UpdateQuestion, q.Reward)
	q.SetTxnChecker(q)
	q.SetInit(q)
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

func (q *Question) InitChain() {
	err := q.asset.AddBalance(common.HexToAddress("0x110e2F71F7a94ba18dbeC96234CC399a2cE61E5D"), big.NewInt(100000))
	if err != nil {
		logrus.Fatal("set balance error: ", err)
	}
}

func (q *Question) AddQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	asker := ctx.GetCaller()
	req := &types.QuestionAddRequest{}
	err := ctx.Bindjson(req)
	if err != nil {
		return err
	}

	err = q.lockForReward(asker, req.TotalRewards)
	if err != nil {
		return err
	}

	scheme := &types.QuestionScheme{
		ID:           ctx.Txn.TxnHash.String(),
		Title:        req.Title,
		Asker:        asker,
		Tags:         req.Tags,
		TotalRewards: req.TotalRewards,
		Timestamp:    req.Timestamp,
		Recommender:  req.Recommender,
	}
	err = q.setQuestion(scheme)
	if err != nil {
		return err
	}
	return ctx.EmitEvent(fmt.Sprintf("add question(%s) successfully by asker(%s)! question-id=%s", scheme.Title, asker.String(), scheme.ID))
}

func (q *Question) UpdateQuestion(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	asker := ctx.GetCaller()
	req := &types.QuestionUpdateRequest{}
	err := ctx.Bindjson(req)
	if err != nil {
		return err
	}

	question, err := q.getQuestion(req.ID)
	if err != nil {
		return err
	}
	if question.Asker != asker {
		return types.ErrNoPermission
	}

	err = q.unlockForReward(asker, question.TotalRewards)
	if err != nil {
		return err
	}
	err = q.lockForReward(asker, req.TotalRewards)
	if err != nil {
		return err
	}

	scheme := &types.QuestionScheme{
		ID:           req.ID,
		Title:        req.Title,
		Asker:        asker,
		Tags:         req.Tags,
		TotalRewards: req.TotalRewards,
		Timestamp:    req.Timestamp,
		Recommender:  req.Recommender,
	}

	err = q.setQuestion(scheme)
	if err != nil {
		return err
	}
	return ctx.EmitEvent(fmt.Sprintf("update question(%s) successfully!", req.ID))
}

func (q *Question) Reward(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	req := &types.RewardRequest{}
	err := ctx.Bindjson(req)
	if err != nil {
		return err
	}

	question, err := q.getQuestion(req.QID)
	if err != nil {
		return err
	}
	for answerID, reward := range req.Rewards {
		answer, err := q.answer.getAnswer(answerID)
		if err != nil {
			return err
		}
		if reward.Cmp(question.TotalRewards) > 0 {
			return types.ErrRewardNotEnough
		}
		err = q.asset.AddBalance(answer.Answerer, reward)
		if err != nil {
			return err
		}
		question.TotalRewards = new(big.Int).Sub(question.TotalRewards, reward)
	}

	return q.setQuestion(question)
}

func (q *Question) setQuestion(scheme *types.QuestionScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	q.State.Set(q, []byte(scheme.ID), byt)
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

func (q *Question) lockForReward(addr common.Address, amount *big.Int) error {
	if amount.Sign() <= 0 {
		return types.ErrRewardIllegal
	}
	balance := q.asset.GetBalance(addr)
	if balance.Cmp(amount) < 0 {
		return errors.New("not enough balance for rewards")
	}
	return q.asset.SubBalance(addr, amount)
}

func (q *Question) unlockForReward(addr common.Address, amount *big.Int) error {
	return q.asset.AddBalance(addr, amount)
}

func checkOffchainOrStoreOnchain(fromP2P bool, info *types.StoreInfo, store filestore.FileStore) error {
	if !fromP2P {
		hash, err := store.Put("", info)
		if err != nil {
			return err
		}
		info.Hash = hash
		info.Url = store.Url()
		info.Content = nil
		return nil
	}
	byt, err := store.Get(info.Hash)
	if err != nil {
		return err
	}
	if byt == nil {
		return types.ErrFileNotfound
	}
	return nil
}
