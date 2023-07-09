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
	ID          string `json:"id" gorm:"primaryKey;column:id"`
	QID         string `json:"qid" gorm:"column:qid"`
	FileHash    string `json:"file_hash" gorm:"column:file_hash"`
	Answerer    string `json:"answerer" gorm:"column:answerer"`
	Timestamp   string `json:"timestamp" gorm:"column:timestamp"`
	Recommender string `json:"recommender" gorm:"column:recommender"`
}
