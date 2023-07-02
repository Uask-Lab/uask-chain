package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"uask-chain/core"
	"uask-chain/db"
	"uask-chain/filestore"
	"uask-chain/search"
)

func main() {
	localStore, err := filestore.NewLocalStore("uask/files")
	if err != nil {
		logrus.Fatal(err)
	}
	meili, err := search.NewMeili(&search.MeiliCfg{
		// this host is for docker
		Host:       "http://meili:7700",
		Index:      "uask",
		PrimaryKey: "id",
	})
	if err != nil {
		logrus.Fatal(err)
	}
	database, err := db.NewDB("uask/db/scheme")
	if err != nil {
		logrus.Fatal(err)
	}

	poaCfg := new(poa.PoaConfig)
	config.LoadTomlConf("poa.toml", poaCfg)

	figure.NewColorFigure("Uask", "big", "green", false).Print()

	startup.InitConfigFromPath("yu.toml")
	startup.DefaultStartup(
		poa.NewPoa(poaCfg),
		core.NewQuestion(localStore, meili, database),
		core.NewAnswer(localStore, database),
		core.NewComment(localStore, database),
	)
}
