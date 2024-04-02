package storage

type Storage interface {
	Pull(id string) ([]byte, error)
	GetValue(id string) ([]byte, error)
	Save(key string, value []byte) error
	Delete(key string) error
}
