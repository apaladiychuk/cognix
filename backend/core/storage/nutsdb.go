package storage

import (
	"cognix.ch/api/v2/core/utils"
	"github.com/nutsdb/nutsdb"
)

const bucket = "internal-storage"

type nutsDbStorage struct {
	db *nutsdb.DB
}

func NewNutsDbStorage(dbPath string) (Storage, error) {
	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(dbPath),
	)
	if err != nil {
		return nil, err
	}
	if err = db.Update(func(tx *nutsdb.Tx) error {
		if !tx.ExistBucket(nutsdb.DataStructureBTree, bucket) {
			return tx.NewBucket(nutsdb.DataStructureBTree, bucket)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &nutsDbStorage{db: db}, nil
}

func (s *nutsDbStorage) GetValue(id string) ([]byte, error) {
	var value []byte
	var err error
	if err = s.db.View(func(tx *nutsdb.Tx) error {
		value, err = tx.Get(bucket, []byte(id))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, utils.NotFound.Wrap(err, "value does not exists")
	}
	return value, nil
}

func (s *nutsDbStorage) Save(key string, value []byte) error {
	if err := s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucket, []byte(key), value, nutsdb.Persistent)
	}); err != nil {
		return utils.Internal.Wrap(err, "can not save value")
	}
	return nil
}

func (s *nutsDbStorage) Delete(key string) error {
	if err := s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(bucket, []byte(key))
	}); err != nil {
		return utils.Internal.Wrap(err, "can not delete value")
	}
	return nil
}
