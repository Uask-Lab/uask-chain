package question

import (
	"github.com/yu-org/yu/core/context"
	"uask-chain/types"
)

const (
	VoteQuestionReputationNeed = 3
	UpVoteQuestionReputation   = 2
	DownVoteQuestionReputation = 2
)

func (q *Question) UpVote(ctx *context.WriteContext) error {
	upVoter := ctx.GetCaller()
	err := q.user.CheckReputation(upVoter, VoteQuestionReputationNeed)
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

	return q.user.IncreaseReputation(upVoter, UpVoteQuestionReputation)
}

func (q *Question) DownVote(ctx *context.WriteContext) error {
	downVoter := ctx.GetCaller()
	err := q.user.CheckReputation(downVoter, VoteQuestionReputationNeed)
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

	return q.user.IncreaseReputation(downVoter, DownVoteQuestionReputation)
}
