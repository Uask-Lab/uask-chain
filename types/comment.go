package types

import "github.com/yu-org/yu/common"

type CommentInfo struct {
	CommentUpdateRequest
	Commenter common.Address `json:"commenter"`
}

type CommentAddRequest struct {
	// reply question id
	QID string `json:"qid"`
	// reply answer id
	AID string `json:"aid"`
	// reply comment id
	CID       string `json:"cid"`
	Content   []byte `json:"content"`
	Timestamp string `json:"timestamp"`
}

type CommentUpdateRequest struct {
	ID string `json:"id"`
	CommentAddRequest
}

// CommentScheme stores into statedb
type CommentScheme struct {
	ID string `json:"id" gorm:"primaryKey"`
	// reply question id
	QID string `json:"qid"`
	// reply answer id
	AID string `json:"aid"`
	// reply comment id
	CID       string `json:"cid"`
	FileHash  string `json:"file_hash"`
	Commenter string `json:"commenter"`
	Timestamp string `json:"timestamp"`
}
