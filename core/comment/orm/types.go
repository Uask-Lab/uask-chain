package orm

// CommentScheme stores into statedb
type CommentScheme struct {
	ID string `json:"id" gorm:"primaryKey;column:id"`
	// reply question id
	QID string `json:"qid" gorm:"column:qid"`
	// reply answer id
	AID       string `json:"aid" gorm:"column:aid"`
	Content   string `json:"content" gorm:"column:content"`
	Commenter string `json:"commenter" gorm:"column:commenter"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
}

func (CommentScheme) TableName() string {
	return "comment"
}
