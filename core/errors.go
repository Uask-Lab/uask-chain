package core

import "github.com/pkg/errors"

var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotFound   = errors.New("answer not found")
	ErrCommentNotFound  = errors.New("comment not found")
	ErrNoPermission     = errors.New("no permission")
)
