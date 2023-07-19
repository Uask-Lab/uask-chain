package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	ycfg "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"uask-chain/config"
	"uask-chain/core"
	"uask-chain/db"
	"uask-chain/filestore"
	"uask-chain/search"
)

func main() {
	uaskCfg := new(config.Config)
	ycfg.LoadTomlConf("./cfg/uask.toml", uaskCfg)

	localStore, err := filestore.NewLocalStore(uaskCfg.Files)
	if err != nil {
		logrus.Fatal(err)
	}
	meili, err := search.NewMeili(uaskCfg.Search)
	if err != nil {
		logrus.Fatal(err)
	}
	database, err := db.NewDB(uaskCfg.DB)
	if err != nil {
		logrus.Fatal(err)
	}

	poaCfg := new(poa.PoaConfig)
	ycfg.LoadTomlConf("./cfg/poa.toml", poaCfg)

	figure.NewColorFigure("Uask", "big", "green", false).Print()

	startup.InitConfigFromPath("./cfg/yu.toml")
	startup.DefaultStartup(
		poa.NewPoa(poaCfg),
		core.NewQuestion(localStore, meili, database),
		core.NewAnswer(localStore, database),
		core.NewComment(localStore, database),
	)
}
