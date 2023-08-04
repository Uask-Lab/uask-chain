package types

type AnswerInfo struct {
	ID        string `json:"id"`
	QID       string `json:"qid"`
	Content   string `json:"content"`
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

// AnswerScheme stores into statedb
type AnswerScheme struct {
	ID        string `json:"id" gorm:"primaryKey;column:id"`
	QID       string `json:"qid" gorm:"column:qid"`
	FileHash  string `json:"file_hash" gorm:"column:file_hash"`
	Answerer  string `json:"answerer" gorm:"column:answerer"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
}
