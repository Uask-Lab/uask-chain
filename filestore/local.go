package filestore

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func (l *LocalStore) Put(key string, file []byte) (string, error) {
	f, err := os.Open(filepath.Join(l.dir, key))
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(file)
	return key, err
}

func (l *LocalStore) Get(key string) ([]byte, error) {
	f, err := os.Open(filepath.Join(l.dir, key))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
