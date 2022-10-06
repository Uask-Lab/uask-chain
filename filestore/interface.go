package filestore

import "uask-chain/types"

type FileStore interface {
	Put(key string, content *types.StoreInfo) (hash string, err error)
	Get(hash string) ([]byte, error)
	Exist(hash string) bool
}
