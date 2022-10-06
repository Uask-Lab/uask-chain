package types

import "github.com/yu-org/yu/common"

type AnswerInfo struct {
	AnswerUpdateRequest
	Answerer common.Address `json:"answerer"`
	Comments []*CommentInfo `json:"comments"`
}

type AnswerAddRequest struct {
	// question id
	QID         string         `json:"qid"`
	Content     *StoreInfo     `json:"content"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type AnswerUpdateRequest struct {
	ID string `json:"id"`
	AnswerAddRequest
}

type AnswerScheme struct {
	ID          string         `json:"id"`
	QID         string         `json:"qid"`
	Answerer    common.Address `json:"answerer"`
	ContentStub string         `json:"content_stub"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}
