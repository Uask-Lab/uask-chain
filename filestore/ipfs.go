package filestore

import (
	api "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"os"
	"path/filepath"
	"uask-chain/types"
)

type IpfsStore struct {
	cli *api.Shell
	dir string
}

func NewIpfsStore(url, dir string) (*IpfsStore, error) {
	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return nil, err
	}
	cli := api.NewShell(url)
	return &IpfsStore{cli: cli, dir: dir}, nil
}

func (i *IpfsStore) Put(_ string, content *types.StoreInfo) (string, error) {
	panic("implement me")
	// return i.cli.Add(bytes.NewReader(content.Content))
}

func (i *IpfsStore) Get(hash string) ([]byte, error) {
	err := i.cli.Get(hash, i.dir)
	if err != nil {
		return nil, err
	}
	fpath := filepath.Join(i.dir, hash)
	byt, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	_ = os.RemoveAll(fpath)
	return byt, nil
}

func (i *IpfsStore) Exist(hash string) bool {
	//TODO implement me
	panic("implement me")
}
