package types

import (
	"github.com/yu-org/yu/common"
	"math/big"
)

type QuestionInfo struct {
	QuestionUpdateRequest
	Asker   common.Address `json:"asker"`
	Answers []*AnswerInfo  `json:"answers"`
}

type QuestionAddRequest struct {
	Title       string         `json:"title"`
	Content     []byte         `json:"content"`
	Tags        []string       `json:"tags"`
	Reward      *big.Int       `json:"reward"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type QuestionUpdateRequest struct {
	ID string `json:"id"`
	QuestionAddRequest
}

type QuestionScheme struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Asker       common.Address `json:"asker"`
	ContentStub string         `json:"content_stub"`
	Tags        []string       `json:"tags"`
	Reward      *big.Int       `json:"reward"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}
