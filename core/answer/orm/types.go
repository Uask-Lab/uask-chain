package orm

// AnswerScheme stores into statedb
type AnswerScheme struct {
	ID        string `json:"id" gorm:"primaryKey;column:id"`
	QID       string `json:"qid" gorm:"column:qid"`
	FileHash  string `json:"file_hash" gorm:"column:file_hash"`
	Answerer  string `json:"answerer" gorm:"column:answerer"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
}

func (AnswerScheme) TableName() string {
	return "answer"
}
