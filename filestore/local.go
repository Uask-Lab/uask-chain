package filestore

import (
	"io"
	"os"
	"path/filepath"
	"uask-chain/types"
)

type LocalStore struct {
	dir string
}

func NewLocalStore(dir string) (*LocalStore, error) {
	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return nil, err
	}
	return &LocalStore{dir: dir}, err
}

func (l *LocalStore) Put(key string, file *types.StoreInfo) (string, error) {
	f, err := os.Open(filepath.Join(l.dir, key))
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(file.Content)
	return key, err
}

func (l *LocalStore) Get(key string) ([]byte, error) {
	f, err := os.Open(filepath.Join(l.dir, key))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func (l *LocalStore) Url() string {
	return ""
}

func (l *LocalStore) Exist(key string) bool {
	path := filepath.Join(l.dir, key)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
