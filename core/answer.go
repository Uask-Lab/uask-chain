package core

import (
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"uask-chain/filestore"
)

type Answer struct {
	*tripod.Tripod
	fileStore filestore.FileStore
}

func NewAnswer(fileStore filestore.FileStore) *Answer {
	tri := tripod.NewTripod("answer")
	a := &Answer{Tripod: tri, fileStore: fileStore}
	a.SetExec(a.AddAnswer)
	return a
}

func (a *Answer) AddAnswer(ctx *context.Context) error {

}
