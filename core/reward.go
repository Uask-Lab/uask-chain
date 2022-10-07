package core

import (
	"github.com/pkg/errors"
	"github.com/yu-org/yu/apps/asset"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"math/big"
	"uask-chain/types"
)

type Reward struct {
	*tripod.Tripod
	asset *asset.Asset `tripod:"asset"`
}

func NewReward() *Reward {
	tri := tripod.NewTripod("reward")
	r := &Reward{Tripod: tri}
	r.SetExec(r.Reward)
	return r
}

func (r *Reward) Reward(ctx *context.Context) error {
	ctx.SetLei(10)
	req := ctx.ParamsValue.(*types.RewardRequest)
	for answerID, reward := range req.Rewards {

	}
}

func (r *Reward) LockForReward(addr common.Address, amount *big.Int) error {
	balance := r.asset.GetBalance(addr)
	if balance.Cmp(amount) < 0 {
		return errors.New("not enough balance for rewards")
	}
	return r.asset.SubBalance(addr, amount)
}
