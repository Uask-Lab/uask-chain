package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"uask-chain/core"
	"uask-chain/filestore"
	"uask-chain/search"
)

func main() {
	localStore, err := filestore.NewIpfsStore("localhost:5001", "uask-files")
	if err != nil {
		logrus.Fatal(err)
	}
	meili, err := search.NewMeili(&search.MeiliCfg{
		Host:       "http://localhost:7700",
		Index:      "uask",
		PrimaryKey: "id",
	})
	if err != nil {
		logrus.Fatal(err)
	}

	poaCfg := new(poa.PoaConfig)
	config.LoadTomlConf("poa.toml", poaCfg)

	figure.NewColorFigure("Uask", "big", "green", false).Print()

	startup.InitConfigFromPath("yu.toml")
	startup.DefaultStartup(
		poa.NewPoa(poaCfg),
		core.NewQuestion(localStore, meili),
		core.NewAnswer(localStore),
		core.NewComment(localStore),
	)
}
