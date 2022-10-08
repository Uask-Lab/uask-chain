package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/core/startup"
	"uask-chain/core"
	"uask-chain/filestore"
)

func main() {
	localStore, err := filestore.NewLocalStore("./uask-files")
	if err != nil {
		logrus.Fatal(err)
	}
	startup.StartUpFullNode(
		poa.NewPoa(),
		core.NewQuestion(localStore),
		core.NewAnswer(localStore),
		core.NewComment(localStore),
	)
}
