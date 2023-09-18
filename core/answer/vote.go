package answer

import (
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"uask-chain/types"
)

const (
	VoteAnswerReputationNeed = 5
	PickUpAnswerReputation   = 2
	UpVoteAnswerReputation   = 2
	DownVoteAnswerReputation = 2
)

func (a *Answer) UpVote(ctx *context.WriteContext) error {
	upVoter := ctx.GetCaller()
	err := a.user.CheckReputation(upVoter, VoteAnswerReputationNeed)
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

	return a.user.IncreaseReputation(upVoter, UpVoteAnswerReputation)
}

func (a *Answer) DownVote(ctx *context.WriteContext) error {
	downVoter := ctx.GetCaller()
	err := a.user.CheckReputation(downVoter, VoteAnswerReputationNeed)
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

	err = a.user.ReduceReputation(downVoter, DownVoteAnswerReputation/2)
	if err != nil {
		return err
	}
	return a.user.ReduceReputation(common.HexToAddress(as.Answerer), DownVoteAnswerReputation)
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
	if qs.Asker == asker.String() {
		return types.ErrNoPermission
	}
	err = a.db.PickUp(aid)
	if err != nil {
		return err
	}
	err = a.user.IncreaseReputation(common.HexToAddress(as.Answerer), PickUpAnswerReputation)
	if err != nil {
		return err
	}
	return a.user.IncreaseReputation(asker, PickUpAnswerReputation/2)
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
	err = a.user.ReduceReputation(common.HexToAddress(as.Answerer), PickUpAnswerReputation)
	if err != nil {
		return err
	}
	return a.user.ReduceReputation(asker, PickUpAnswerReputation/2)
}
