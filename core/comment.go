package core

import (
	"encoding/json"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/filestore"
	"uask-chain/types"
)

type Comment struct {
	*tripod.Tripod
	fileStore filestore.FileStore
	// sch       search.Search

	Answer *Answer `tripod:"answer"`
}

func NewComment(fileStore filestore.FileStore) *Comment {
	tri := tripod.NewTripod()
	c := &Comment{Tripod: tri, fileStore: fileStore}
	c.SetWritings(c.AddComment, c.UpdateComment, c.DeleteComment)
	return c
}

func (c *Comment) AddComment(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	commenter := ctx.GetCaller()
	req := &types.CommentAddRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	err = c.ifReplyExist(req.AID, req.CID)
	if err != nil {
		return err
	}

	fileHash, err := c.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:        ctx.Txn.TxnHash.String(),
		AID:       req.AID,
		CID:       req.CID,
		FileHash:  fileHash,
		Commenter: commenter,
		Timestamp: req.Timestamp,
	}
	err = c.setCommentScheme(scheme)
	if err != nil {
		return err
	}
	// add search
	//contentByt, err := c.fileStore.Get(req.Content.Hash)
	//if err != nil {
	//	return err
	//}
	//err = c.sch.AddDoc(&types.Comment{
	//	ID:          scheme.ID,
	//	FileContent: contentByt,
	//	Commenter:   scheme.Commenter,
	//	Timestamp:   scheme.Timestamp,
	//})
	//if err != nil {
	//	return err
	//}

	ctx.EmitStringEvent("add comment(%s) successfully by commenter(%s)", scheme.ID, commenter.String())
	return nil
}

func (c *Comment) UpdateComment(ctx *context.WriteContext) error {
	ctx.SetLei(10)

	commenter := ctx.GetCaller()
	req := &types.CommentUpdateRequest{}
	err := ctx.BindJson(req)
	if err != nil {
		return err
	}

	comment, err := c.getCommentScheme(req.ID)
	if err != nil {
		return err
	}
	if comment.Commenter != commenter {
		return types.ErrNoPermission
	}

	err = c.ifReplyExist(req.AID, req.CID)
	if err != nil {
		return err
	}

	// remove old answer and store new one.
	err = c.fileStore.Remove(comment.FileHash)
	if err != nil {
		return err
	}
	fileHash, err := c.fileStore.Put(req.Content)
	if err != nil {
		return err
	}

	scheme := &types.CommentScheme{
		ID:        req.ID,
		AID:       req.AID,
		CID:       req.CID,
		FileHash:  fileHash,
		Commenter: commenter,
		Timestamp: req.Timestamp,
	}
	err = c.setCommentScheme(scheme)
	if err != nil {
		return err
	}
	// update search
	//contentByt, err := c.fileStore.Get(req.Content.Hash)
	//if err != nil {
	//	return err
	//}
	//err = c.sch.UpdateDoc(scheme.ID, &types.Comment{
	//	ID:          scheme.ID,
	//	FileContent: contentByt,
	//	Commenter:   scheme.Commenter,
	//	Timestamp:   scheme.Timestamp,
	//})
	//if err != nil {
	//	return err
	//}

	ctx.EmitStringEvent("update comment(%s) successfully!", req.ID)
	return nil
}

func (c *Comment) DeleteComment(ctx *context.WriteContext) error {
	ctx.SetLei(10)
	id := ctx.GetString("id")
	commenter := ctx.GetCaller()
	scheme, err := c.getCommentScheme(id)
	if err != nil {
		return err
	}
	if commenter != scheme.Commenter {
		return types.ErrNoPermission
	}
	c.Delete([]byte(id))
	// return c.sch.DeleteDoc(id)
	return nil
}

func (c *Comment) setCommentScheme(scheme *types.CommentScheme) error {
	byt, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	c.Set([]byte(scheme.ID), byt)
	return nil
}

func (c *Comment) existComment(id string) bool {
	return c.Exist([]byte(id))
}

func (c *Comment) getCommentScheme(id string) (*types.CommentScheme, error) {
	byt, err := c.Get([]byte(id))
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
		if !c.Answer.existAnswer(answerID) {
			return types.ErrAnswerNotFound
		}
	}
	if commentID != "" {
		if !c.existComment(commentID) {
			return types.ErrCommentNotFound
		}
	}
	return nil
}
