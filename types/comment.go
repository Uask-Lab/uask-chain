package types

type CommentInfo struct {
	ID string `json:"id"`
	// reply question id
	QID string `json:"qid"`
	// reply answer id
	AID       string `json:"aid"`
	Content   string `json:"content"`
	Commenter string `json:"commenter"`
	Timestamp int64  `json:"timestamp"`
}

type CommentAddRequest struct {
	// reply question id
	QID string `json:"qid"`
	// reply answer id
	AID     string `json:"aid"`
	Content string `json:"content"`
	// Timestamp int64  `json:"timestamp"`
}

type CommentUpdateRequest struct {
	ID string `json:"id"`
	CommentAddRequest
}
