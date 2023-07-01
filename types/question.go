package types

import (
	"github.com/yu-org/yu/common"
)

type QuestionInfo struct {
	QuestionDoc
	AnswersIDs []string `json:"answers_ids"`
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
	ID          string         `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title"`
	FileHash    string         `json:"file_hash"`
	Asker       common.Address `json:"asker"`
	Tags        []string       `json:"tags,omitempty"`
	Timestamp   string         `json:"timestamp"`
	Recommender common.Address `json:"recommender"`
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
