package core

import (
	"encoding/json"
	"fmt"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/core/types"
	"uask-chain/filestore"
)

type Comment struct {
	*tripod.Tripod
	fileStore filestore.FileStore
}

func NewComment(fileStore filestore.FileStore) *Comment {
	tri := tripod.NewTripod("comment")
	c := &Comment{Tripod: tri, fileStore: fileStore}
	c.SetExec(c.AddComment).SetExec(c.UpdateComment)
	return c
}

func (c *Comment) AddComment(ctx *context.Context) error {
	ctx.SetLei(10)

	commenter := ctx.Caller
	req := &types.CommentAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	err = c.ifReplyExist(req.AID, req.CID)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s%s%s%s", commenter.String(), req.AID, req.CID, req.Timestamp)
	stub, err := c.fileStore.Put(id, req.Content)
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:          id,
		AID:         req.AID,
		CID:         req.CID,
		Commenter:   commenter,
		ContentStub: stub,
		Timestamp:   req.Timestamp,
	}
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}
	c.State.Set(c, []byte(id), byt)
	ctx.EmitEvent(fmt.Sprintf("add comment(%s) successfully by commenter(%s)", scheme.ID, commenter.String()))
	return nil
}

func (c *Comment) UpdateComment(ctx *context.Context) error {
	ctx.SetLei(10)

	commenter := ctx.Caller
	req := &types.CommentUpdateRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	comment, err := c.getComment(req.ID)
	if err != nil {
		return err
	}
	if comment.Commenter != commenter {
		return ErrNoPermission
	}

	err = c.ifReplyExist(req.AID, req.CID)
	if err != nil {
		return err
	}

	stub, err := c.fileStore.Put(req.ID, req.Content)
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:          req.ID,
		AID:         req.AID,
		CID:         req.CID,
		Commenter:   commenter,
		ContentStub: stub,
		Timestamp:   req.Timestamp,
	}
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	c.State.Set(c, []byte(req.ID), byt)
	ctx.EmitEvent(fmt.Sprintf("update comment(%s) successfully!", req.ID))
	return nil
}

func (c *Comment) existComment(id string) bool {
	return c.State.Exist(c, []byte(id))
}

func (c *Comment) getComment(id string) (*types.CommentScheme, error) {
	byt, err := c.State.Get(c, []byte(id))
	if err != nil {
		return nil, err
	}
	scheme := &types.CommentScheme{}
	err = json.Unmarshal(byt, scheme)
	if err != nil {
		return nil, err
	}
	return scheme, nil
}

func (c *Comment) ifReplyExist(answerID, commentID string) error {
	if answerID != "" {
		a := c.GetTripod("answer").(*Answer)
		if !a.existAnswer(answerID) {
			return ErrAnswerNotFound
		}
	}
	if commentID != "" {
		if !c.existComment(commentID) {
			return ErrCommentNotFound
		}
	}
	return nil
}
