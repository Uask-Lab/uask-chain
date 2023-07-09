package types

import (
	"github.com/yu-org/yu/common"
)

type QuestionInfo struct {
	QuestionDoc
}

type QuestionAddRequest struct {
	Title       string         `json:"title"`
	Content     []byte         `json:"content"`
	Tags        []string       `json:"tags,omitempty"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}

type QuestionUpdateRequest struct {
	ID string `json:"id"`
	QuestionAddRequest
}

type QuestionScheme struct {
	ID          string   `json:"id" gorm:"primaryKey;column:id"`
	Title       string   `json:"title" gorm:"column:title"`
	FileHash    string   `json:"file_hash" gorm:"column:file_hash"`
	Asker       string   `json:"asker" gorm:"column:asker"`
	Tags        []string `json:"tags,omitempty" gorm:"type:text[];column:tags"`
	Timestamp   string   `json:"timestamp" gorm:"column:timestamp"`
	Recommender string   `json:"recommender" gorm:"column:recommender"`
}

// QuestionDoc stores into search
type QuestionDoc struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Content     []byte         `json:"content"`
	Asker       common.Address `json:"asker"`
	Tags        []string       `json:"tags,omitempty"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
}
