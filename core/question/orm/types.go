package orm

type QuestionScheme struct {
	ID        string   `json:"id" gorm:"primaryKey;column:id"`
	Title     string   `json:"title" gorm:"column:title"`
	FileHash  string   `json:"file_hash" gorm:"column:file_hash"`
	Asker     string   `json:"asker" gorm:"column:asker"`
	UpVotes   uint64   `json:"up_votes" gorm:"column:up_votes"`
	DownVotes uint64   `json:"down_votes" gorm:"column:down_votes"`
	Tags      []string `json:"tags,omitempty" gorm:"type:text[];column:tags"`
	Timestamp int64    `json:"timestamp" gorm:"column:timestamp"`
}

func (QuestionScheme) TableName() string {
	return "question"
}
