package filestore

type FileStore interface {
	Put(key string, file []byte) (stub string, err error)
	Get(stub string) ([]byte, error)
}
