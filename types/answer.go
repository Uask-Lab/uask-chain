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

// AnswerScheme stores into statedb
type AnswerScheme struct {
	ID          string         `json:"id"`
	QID         string         `json:"qid"`
	FileHash    string         `json:"file_hash"`
	Answerer    common.Address `json:"answerer"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

// Answer stores into search
type Answer struct {
	ID          string         `json:"id"`
	Answerer    common.Address `json:"answerer"`
	FileContent []byte         `json:"file_content"`
	Recommender common.Address `json:"recommender"`
	Timestamp   string         `json:"timestamp"`
}
