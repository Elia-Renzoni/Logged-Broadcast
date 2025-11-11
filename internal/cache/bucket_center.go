package cache

import (
	badger "github.com/dgraph-io/badger"
	"bytes"
)


type MemoryCache interface {
	OpenDB()
	SetBucket(key, value string)
	DeleteBucket(key string) error
	FetchBucket(key string) string
}

type Bcache struct {
	db *badger.DB
	err error
	opts badger.Options
}

func (b *Bcache) OpenDB() {
	b.opts = badger.DefaultOptions("./cache")
	b.db, b.err = badger.Open(b.opts)
	if b.err != nil {
		return
	}
}

func (b *Bcache) CloseDB() {
	b.db.Close()
}

func (b Bcache) SetBucket(key, value string) {
	err := b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	if err != nil {
		return
	}
}

func (b Bcache) DeleteBucket(key string) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})

	return err
}

func (b Bcache) FetchBucket(key string) string {
	var value []byte
	var builder bytes.Buffer
	
	err := b.db.View(func(txn *badger.Txn) error {
		val, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		var copyErr error
		value, copyErr = val.ValueCopy(nil)
		return copyErr
	})

	if err != nil {
		return ""
	}

	builder.Write(value)
	return builder.String()
}

