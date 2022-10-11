package types

import "github.com/pkg/errors"

var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotFound   = errors.New("answer not found")
	ErrCommentNotFound  = errors.New("comment not found")

	ErrNoPermission = errors.New("no permission")

	ErrRewardNotEnough = errors.New("reward not enough")
	ErrRewardIllegal   = errors.New("reward is illegal")

	ErrFileNotMatchHash = errors.New("file-content not match file-hash")
)