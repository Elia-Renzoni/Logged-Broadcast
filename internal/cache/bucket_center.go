package cache

import (
	"bytes"
	"log"

	badger "github.com/dgraph-io/badger"
)

type MemoryCache interface {
	OpenDB()
	SetBucket(key, value string)
	DeleteBucket(key string) error
	FetchBucket(key string) string
	CloseDB()
}

type Bcache struct {
	db   *badger.DB
	err  error
	opts badger.Options
}

func (b *Bcache) OpenDB() {
	b.opts = badger.DefaultOptions("./cache")
	b.db, b.err = badger.Open(b.opts)
	if b.err != nil {
		log.Fatalf("%s", b.err.Error())
	}
}

func (b *Bcache) CloseDB() {
	err := b.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (b *Bcache) SetBucket(key, value string) {
	err := b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func (b *Bcache) DeleteBucket(key string) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})

	return err
}

func (b *Bcache) FetchBucket(key string) string {
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
