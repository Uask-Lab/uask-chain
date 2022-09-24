package core

import (
	"github.com/yu-org/yu/common"
	"math/big"
)

type QuestionInfo struct {
	ID string `json:"id"`
	QuestionRequest
	Asker   common.Address `json:"asker"`
	Answers []*AnswerInfo  `json:"answers"`
}

type QuestionRequest struct {
	Title       string         `json:"title"`
	Content     []byte         `json:"content"`
	Tags        []string       `json:"tags"`
	Reward      *big.Int       `json:"reward"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type QuestionScheme struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	ContentStub string         `json:"content_stub"`
	Tags        []string       `json:"tags"`
	Reward      *big.Int       `json:"reward"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type AnswerInfo struct {
	ID string `json:"id"`
	AnswerRequest
	Answerer common.Address `json:"answerer"`
	Comments []*CommentInfo `json:"comments"`
}

type AnswerRequest struct {
	// question id
	QID         string         `json:"qid"`
	Content     []byte         `json:"content"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type AnswerScheme struct {
	ID          string         `json:"id"`
	QID         string         `json:"qid"`
	ContentStub string         `json:"content_stub"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type CommentInfo struct {
	ID string `json:"id"`
	CommentRequest
	Commenter common.Address `json:"commenter"`
}

type CommentRequest struct {
	// reply answer or comment
	ReplyID   string `json:"reply_id"`
	Content   []byte `json:"content"`
	Timestamp string `json:"timestamp"`
}

type CommentScheme struct {
	ID          string `json:"id"`
	ReplyID     string `json:"reply_id"`
	ContentStub string `json:"content_stub"`
}
