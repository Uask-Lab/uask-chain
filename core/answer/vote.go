package answer

import (
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"uask-chain/types"
)

func (a *Answer) UpVote(ctx *context.WriteContext) error {
	upVoter := ctx.GetCaller()
	err := a.user.CheckReputation(upVoter, types.VoteAnswerReputationNeed)
	if err != nil {
		return err
	}
	aid := ctx.GetString("id")
	as, err := a.db.GetAnswer(aid)
	if err != nil {
		return err
	}

	if as.Answerer == upVoter.String() {
		return types.ErrCannotVoteYourself
	}

	err = a.db.UpVote(aid)
	if err != nil {
		return err
	}

	err = a.user.IncreaseReputation(upVoter, types.UpVoteAnswerReputationIncrease)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"status": "success"})
}

func (a *Answer) DownVote(ctx *context.WriteContext) error {
	downVoter := ctx.GetCaller()
	err := a.user.CheckReputation(downVoter, types.VoteAnswerReputationNeed)
	if err != nil {
		return err
	}
	aid := ctx.GetString("id")
	as, err := a.db.GetAnswer(aid)
	if err != nil {
		return err
	}

	if as.Answerer == downVoter.String() {
		return types.ErrCannotVoteYourself
	}

	err = a.db.DownVote(aid)
	if err != nil {
		return err
	}

	err = a.user.ReduceReputation(downVoter, types.DownVoteAnswerReputationReduce/2)
	if err != nil {
		return err
	}
	err = a.user.ReduceReputation(common.HexToAddress(as.Answerer), types.DownVoteAnswerReputationReduce)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"status": "success"})
}

func (a *Answer) PickUp(ctx *context.WriteContext) error {
	asker := ctx.GetCaller()
	aid := ctx.GetString("id")
	as, err := a.db.GetAnswer(aid)
	if err != nil {
		return err
	}
	qs, err := a.Question.GetQ(as.QID)
	if err != nil {
		return err
	}
	if qs.Asker != asker.String() {
		return types.ErrNoPermission
	}
	err = a.db.PickUp(aid)
	if err != nil {
		return err
	}
	err = a.user.IncreaseReputation(common.HexToAddress(as.Answerer), types.PickUpAnswerReputationIncrease)
	if err != nil {
		return err
	}
	err = a.user.IncreaseReputation(asker, types.PickUpAnswerReputationIncrease/2)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"status": "success"})
}

func (a *Answer) Drop(ctx *context.WriteContext) error {
	asker := ctx.GetCaller()
	aid := ctx.GetString("id")
	as, err := a.db.GetAnswer(aid)
	if err != nil {
		return err
	}
	qs, err := a.Question.GetQ(as.QID)
	if err != nil {
		return err
	}
	if qs.Asker == asker.String() {
		return types.ErrNoPermission
	}
	err = a.db.Drop(aid)
	if err != nil {
		return err
	}
	err = a.user.ReduceReputation(common.HexToAddress(as.Answerer), types.PickUpAnswerReputationIncrease)
	if err != nil {
		return err
	}
	err = a.user.ReduceReputation(asker, types.PickUpAnswerReputationIncrease/2)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"status": "success"})
}
