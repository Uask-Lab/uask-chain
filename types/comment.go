package types

import "github.com/yu-org/yu/common"

type CommentInfo struct {
	CommentUpdateRequest
}

type CommentAddRequest struct {
	// reply question id
	QID string `json:"qid"`
	// reply answer id
	AID string `json:"aid"`
	// reply comment id
	CID       string         `json:"cid"`
	Content   []byte         `json:"content"`
	Timestamp string         `json:"timestamp"`
	Commenter common.Address `json:"commenter"`
}

type CommentUpdateRequest struct {
	ID string `json:"id"`
	CommentAddRequest
}

// CommentScheme stores into statedb
type CommentScheme struct {
	ID string `json:"id" gorm:"primaryKey;column:id"`
	// reply question id
	QID string `json:"qid" gorm:"column:qid"`
	// reply answer id
	AID string `json:"aid" gorm:"column:aid"`
	// reply comment id
	CID       string `json:"cid" gorm:"column:cid"`
	FileHash  string `json:"file_hash" gorm:"column:file_hash"`
	Commenter string `json:"commenter" gorm:"column:commenter"`
	Timestamp string `json:"timestamp" gorm:"column:timestamp"`
}
