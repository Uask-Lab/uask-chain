package types

type AnswerInfo struct {
	ID        string `json:"id"`
	QID       string `json:"qid"`
	Content   string `json:"content"`
	Answerer  string `json:"answerer"`
	Timestamp int64  `json:"timestamp"`
}

type AnswerAddRequest struct {
	// question id
	QID     string `json:"qid"`
	Content string `json:"content"`
}

type AnswerUpdateRequest struct {
	ID string `json:"id"`
	AnswerAddRequest
}
