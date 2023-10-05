package types

import (
	"errors"
	"fmt"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotFound   = errors.New("answer not found")
	ErrCommentNotFound  = errors.New("comment not found")

	ErrQuestionTitleTooLong = errors.New("question title too long")
	ErrCommentTooLong       = errors.New("comment too long")

	ErrNoPermission       = errors.New("no permission")
	ErrCannotVoteYourself = errors.New("you cannot vote yourself")

	ErrNoneToReply = errors.New("none to reply")

	// ErrReputationValueInsufficient = errors.New("reputation value insufficient")
)

type ErrReputationValueInsufficient struct {
	value int64
}

func ReputationValueInsufficientErr(value int64) ErrReputationValueInsufficient {
	return ErrReputationValueInsufficient{value: value}
}

func (rvi ErrReputationValueInsufficient) Error() string {
	return fmt.Sprintf("reputation value insufficient: %d", rvi.value)
}
