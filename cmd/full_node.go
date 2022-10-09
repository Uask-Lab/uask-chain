package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/core/keypair"
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
		poa.NewPoa(&poa.PoaConfig{
			KeyType:  keypair.Sr25519,
			MySecret: "dayu",
			Validators: []*poa.ValidatorConf{
				{Pubkey: "", P2pIp: "12D3KooWHHzSeKaY8xuZVzkLbKFfvNgPPeKhFBGrMbNzbm5akpqu"},
				{Pubkey: "", P2pIp: "12D3KooWSKPs95miv8wzj3fa5HkJ1tH7oEGumsEiD92n2MYwRtQG"},
				{Pubkey: "", P2pIp: "12D3KooWRuwP7nXaRhZrmoFJvPPGat2xPafVmGpQpZs5zKMtwqPH"},
			},
		}),
		core.NewQuestion(localStore),
		core.NewAnswer(localStore),
		core.NewComment(localStore),
	)
}
