package database

import (
	goErrors "errors"
	"fmt"
	"log"
	"os"

	"go.etcd.io/bbolt"
	"go.etcd.io/bbolt/errors"
)

var (
	_               Database = &boltDB{}
	ErrUserNotFound          = goErrors.New("user not found")
)

const (
	bucketName = "Users"
)

type Database interface {
	Read(name string) (string, error)
	Write(name string, age string) error
	Delete(name string) error
	Close() error
}

type boltDB struct {
	DB *bbolt.DB
}

func New(path string, mode os.FileMode, options *bbolt.Options) (*boltDB, error) {
	db, err := bbolt.Open(path, mode, options)
	if err != nil {
		log.Printf("there was an error opening the database: %v", err)
		return nil, err
	}
	return &boltDB{
		DB: db,
	}, nil
}

func (b *boltDB) Read(key string) (string, error) {
	var value string
	err := b.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.ErrBucketNotFound
		}
		val := bucket.Get([]byte(key))
		if val == nil {
			return ErrUserNotFound
		}
		value = string(val)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error reading from the database: %w", err)
	}
	return value, nil
}

func (b *boltDB) Write(key string, val string) error {
	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
		err = bucket.Put([]byte(key), []byte(val))
		if err != nil {
			return fmt.Errorf("error putting key/val into the database: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error writing to the database: %w", err)
	}
	return nil
}

func (b *boltDB) Delete(key string) error {
	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.ErrBucketNotFound
		}
		err := bucket.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("error in deleting key %s from database: %w", key, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting key %s from database: %w", key, err)
	}
	return nil
}

func (b *boltDB) Close() error {
	return b.DB.Close()
}
