package config

import (
	"uask-chain/filestore"
	"uask-chain/search"
)

type Config struct {
	Files     *filestore.Config `toml:"files"`
	Search    *search.MeiliCfg  `toml:"search"`
	DSN       string            `toml:"dsn"`
	WhiteList map[string]uint64 `toml:"white_list"`
}
