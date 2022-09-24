package core

import (
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/filestore"
)

type Comment struct {
	*tripod.Tripod
	fileStore filestore.FileStore
}

func NewComment(fileStore filestore.FileStore) *Comment {
	tri := tripod.NewTripod("comment")
	c := &Comment{Tripod: tri, fileStore: fileStore}
	c.SetExec(c.AddComment)
	return c
}

func (c *Comment) AddComment(ctx *context.Context) error {
	ctx.SetLei(10)

	req := &CommentRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	// TODO: store into file-store
}
