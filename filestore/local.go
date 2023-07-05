package filestore

import (
	"crypto/sha256"
	"github.com/yu-org/yu/common"
	"io"
	"os"
	"path/filepath"
)

type LocalStore struct {
	cfg *Config
}

func NewLocalStore(cfg *Config) (*LocalStore, error) {
	err := os.MkdirAll(cfg.Dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &LocalStore{cfg: cfg}, err
}

func (l *LocalStore) Put(content []byte) (string, error) {
	hashByt := sha256.Sum256(content)
	hash := common.Bytes2Hex(hashByt[:])

	path := filepath.Join(l.cfg.Dir, hash)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(content)
	return hash, err
}

func (l *LocalStore) Get(hash string) ([]byte, error) {
	f, err := os.Open(filepath.Join(l.cfg.Dir, hash))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func (l *LocalStore) Remove(hash string) error {
	return os.RemoveAll(filepath.Join(l.cfg.Dir, hash))
}

func (l *LocalStore) Url() string {
	return ""
}

func (l *LocalStore) Exist(key string) bool {
	path := filepath.Join(l.cfg.Dir, key)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
