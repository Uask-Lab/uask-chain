package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/asset"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"uask-chain/core"
	"uask-chain/filestore"
)

func main() {
	localStore, err := filestore.NewLocalStore("uask-files")
	if err != nil {
		logrus.Fatal(err)
	}

	poaCfg := &poa.PoaConfig{}
	config.LoadTomlConf("poa.toml", poaCfg)

	startup.InitConfigFromPath("yu.toml")
	startup.StartUpFullNode(
		poa.NewPoa(poaCfg),
		asset.NewAsset("uask"),
		core.NewQuestion(localStore),
		core.NewAnswer(localStore),
		core.NewComment(localStore),
	)
}
