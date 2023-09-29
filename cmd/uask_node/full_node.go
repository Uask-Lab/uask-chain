package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	ycfg "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"uask-chain/config"
	"uask-chain/core/answer"
	"uask-chain/core/comment"
	"uask-chain/core/question"
	"uask-chain/core/user"
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
	database, err := gorm.Open(postgres.Open(uaskCfg.DSN), &gorm.Config{CreateBatchSize: 50000})
	if err != nil {
		logrus.Fatal(err)
	}

	poaCfg := new(poa.PoaConfig)
	ycfg.LoadTomlConf("./cfg/poa.toml", poaCfg)

	figure.NewColorFigure("Uask", "big", "green", false).Print()

	startup.InitConfigFromPath("./cfg/yu.toml")
	startup.DefaultStartup(
		poa.NewPoa(poaCfg),
		user.NewUser(database, uaskCfg.WhiteList),
		question.NewQuestion(localStore, meili, database),
		answer.NewAnswer(localStore, database),
		comment.NewComment(localStore, database),
	)
}
