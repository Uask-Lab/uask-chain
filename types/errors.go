package types

import "errors"

var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotFound   = errors.New("answer not found")
	ErrCommentNotFound  = errors.New("comment not found")

	ErrCommentTooLong = errors.New("comment too long")

	ErrNoPermission       = errors.New("no permission")
	ErrCannotVoteYourself = errors.New("you cannot vote yourself")

	ErrNoneToReply = errors.New("none to reply")

	ErrReputationValueInsufficient = errors.New("reputation value insufficient")
)
