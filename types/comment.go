package types

import "github.com/yu-org/yu/common"

type CommentInfo struct {
	CommentUpdateRequest
	Commenter common.Address `json:"commenter"`
}

type CommentAddRequest struct {
	// reply answer id
	AID string `json:"aid"`
	// reply comment id
	CID       string `json:"cid"`
	Content   []byte `json:"content"`
	Timestamp string `json:"timestamp"`
}

type CommentUpdateRequest struct {
	ID string `json:"id"`
	CommentAddRequest
}

type CommentScheme struct {
	ID string `json:"id"`
	// reply answer id
	AID string `json:"aid"`
	// reply comment id
	CID       string         `json:"cid"`
	FileHash  string         `json:"file_hash"`
	Commenter common.Address `json:"commenter"`
	Timestamp string         `json:"timestamp"`
}

// Comment stores into search
type Comment struct {
	ID          string         `json:"id"`
	FileContent []byte         `json:"file_content"`
	Commenter   common.Address `json:"commenter"`
	Timestamp   string         `json:"timestamp"`
}
