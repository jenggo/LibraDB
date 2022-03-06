package main

import (
	"os"
	"sync"
)

type DB struct {
	rwlock   sync.Mutex   // Allows only one writer at a time
	*dal
}

func Open(path string, options *Options) (*DB, error) {
	var err error

	options.pageSize = os.Getpagesize()
	dal, err := newDal(path, options)
	if err != nil {
		return nil, err
	}

	db := &DB{
		sync.Mutex{},
		dal,
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.close()
}

func (db *DB) ReadTx() *tx {
	return newTx(db, db.root, false)
}

func (db *DB) WriteTx() *tx {
	db.rwlock.Lock()
	return newTx(db, db.root, true)
}