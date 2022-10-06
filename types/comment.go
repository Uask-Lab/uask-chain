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
	CID       string     `json:"cid"`
	Content   *StoreInfo `json:"content"`
	Timestamp string     `json:"timestamp"`
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
	CID         string         `json:"cid"`
	Commenter   common.Address `json:"commenter"`
	ContentStub string         `json:"content_stub"`
	Timestamp   string         `json:"timestamp"`
}
