package types

type QuestionInfo struct {
	QuestionDoc
}

type QuestionAddRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
	// Timestamp int64    `json:"timestamp"`
}

type QuestionUpdateRequest struct {
	ID string `json:"id"`
	QuestionAddRequest
}

type QuestionScheme struct {
	ID        string   `json:"id" gorm:"primaryKey;column:id"`
	Title     string   `json:"title" gorm:"column:title"`
	FileHash  string   `json:"file_hash" gorm:"column:file_hash"`
	Asker     string   `json:"asker" gorm:"column:asker"`
	Tags      []string `json:"tags,omitempty" gorm:"type:text[];column:tags"`
	Timestamp int64    `json:"timestamp" gorm:"column:timestamp"`
}

// QuestionDoc stores into search
type QuestionDoc struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Asker     string   `json:"asker"`
	Tags      []string `json:"tags,omitempty"`
	Timestamp int64    `json:"timestamp"`
}
