package storage

type Storage interface {
	GetValue(id string) ([]byte, error)
	Save(key string, value []byte) error
	Delete(key string) error
}
