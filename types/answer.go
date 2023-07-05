package types

import "github.com/yu-org/yu/common"

type AnswerInfo struct {
	AnswerUpdateRequest
	CommentsIDs []string `json:"comments_ids"`
}

type AnswerAddRequest struct {
	// question id
	QID         string         `json:"qid"`
	Content     []byte         `json:"content"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type AnswerUpdateRequest struct {
	ID string `json:"id"`
	AnswerAddRequest
}

// AnswerScheme stores into statedb
type AnswerScheme struct {
	ID          string `json:"id" gorm:"primaryKey"`
	QID         string `json:"qid"`
	FileHash    string `json:"file_hash"`
	Answerer    string `json:"answerer"`
	Timestamp   string `json:"timestamp"`
	Recommender string `json:"recommender"`
}
