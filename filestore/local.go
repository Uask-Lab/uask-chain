package filestore

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
)

type LocalStore struct {
	dir string
}

func NewLocalStore(dir string) (*LocalStore, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &LocalStore{dir: dir}, err
}

func (l *LocalStore) Put(content []byte) (string, error) {
	hashByt := sha256.Sum256(content)
	hash := string(hashByt[:])

	path := filepath.Join(l.dir, hash)
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(content)
	return hash, err
}

func (l *LocalStore) Get(hash string) ([]byte, error) {
	f, err := os.Open(filepath.Join(l.dir, hash))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func (l *LocalStore) Remove(hash string) error {
	return os.RemoveAll(filepath.Join(l.dir, hash))
}

func (l *LocalStore) Url() string {
	return ""
}

func (l *LocalStore) Exist(key string) bool {
	path := filepath.Join(l.dir, key)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
