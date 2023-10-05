package orm

// AnswerScheme stores into statedb
type AnswerScheme struct {
	ID         string `json:"id" gorm:"primaryKey;column:id"`
	QID        string `json:"qid" gorm:"column:qid"`
	FileHash   string `json:"file_hash" gorm:"column:file_hash"`
	Answerer   string `json:"answerer" gorm:"column:answerer"`
	UpVotes    uint64 `json:"up_votes" gorm:"column:up_votes"`
	DownVotes  uint64 `json:"down_votes" gorm:"column:down_votes"`
	IsPickedUp bool   `json:"is_picked_up" gorm:"column:is_picked_up;default:false"`
	Timestamp  int64  `json:"timestamp" gorm:"column:timestamp"`
}

func (AnswerScheme) TableName() string {
	return "answer"
}
