package question

import (
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"uask-chain/types"
)

func (q *Question) UpVote(ctx *context.WriteContext) error {
	upVoter := ctx.GetCaller()
	err := q.user.CheckReputation(upVoter, types.VoteQuestionReputationNeed)
	if err != nil {
		return err
	}
	qid := ctx.GetString("id")
	qs, err := q.db.GetQuestion(qid)
	if err != nil {
		return err
	}
	if qs.Asker == upVoter.String() {
		return types.ErrCannotVoteYourself
	}

	err = q.db.UpVote(qid)
	if err != nil {
		return err
	}

	err = q.user.IncreaseReputation(upVoter, types.UpVoteQuestionReputationIncrease)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"status": "success"})
}

func (q *Question) DownVote(ctx *context.WriteContext) error {
	downVoter := ctx.GetCaller()
	err := q.user.CheckReputation(downVoter, types.VoteQuestionReputationNeed)
	if err != nil {
		return err
	}
	qid := ctx.GetString("id")
	qs, err := q.db.GetQuestion(qid)
	if err != nil {
		return err
	}
	if qs.Asker == downVoter.String() {
		return types.ErrCannotVoteYourself
	}

	err = q.db.DownVote(qid)
	if err != nil {
		return err
	}

	err = q.user.ReduceReputation(downVoter, types.DownVoteQuestionReputationReduce/2)
	if err != nil {
		return err
	}
	err = q.user.ReduceReputation(common.HexToAddress(qs.Asker), types.DownVoteQuestionReputationReduce)
	if err != nil {
		return err
	}
	return ctx.EmitJsonEvent(map[string]string{"status": "success"})
}
