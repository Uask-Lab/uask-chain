package config

import (
	"uask-chain/db"
	"uask-chain/filestore"
	"uask-chain/search"
)

type Config struct {
	Files  *filestore.Config `toml:"files"`
	Search *search.MeiliCfg  `toml:"search"`
	DB     *db.Config        `toml:"db"`
}
