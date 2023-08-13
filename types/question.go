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

// QuestionDoc stores into search
type QuestionDoc struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Asker     string   `json:"asker"`
	Tags      []string `json:"tags,omitempty"`
	Timestamp int64    `json:"timestamp"`
}
