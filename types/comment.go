package types

type CommentInfo struct {
	CommentUpdateRequest
}

type CommentAddRequest struct {
	// reply question id
	QID string `json:"qid"`
	// reply answer id
	AID       string `json:"aid"`
	Content   []byte `json:"content"`
	Timestamp int64  `json:"timestamp"`
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
	AID       string `json:"aid" gorm:"column:aid"`
	FileHash  string `json:"file_hash" gorm:"column:file_hash"`
	Commenter string `json:"commenter" gorm:"column:commenter"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
}
