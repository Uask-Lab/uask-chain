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
	Title        string         `json:"title"`
	Content      *StoreInfo     `json:"content"`
	Tags         []string       `json:"tags"`
	TotalRewards *big.Int       `json:"total_rewards"`
	Timestamp    string         `json:"timestamp"`
	Recommender  common.Address `json:"recommender"`
}

type QuestionUpdateRequest struct {
	ID string `json:"id"`
	QuestionAddRequest
}

type QuestionScheme struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Asker        common.Address `json:"asker"`
	Tags         []string       `json:"tags"`
	TotalRewards *big.Int       `json:"total_rewards"`
	Timestamp    string         `json:"timestamp"`
	Recommender  common.Address `json:"recommender"`
}
