package filestore

type FileStore interface {
	Put(content []byte) (hash string, err error)
	Get(hash string) ([]byte, error)
	Remove(hash string) error
	Url() string
	Exist(hash string) bool
}
