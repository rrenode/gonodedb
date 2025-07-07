package db

import (
	"github.com/dgraph-io/badger/v4"
)

func Open(path string, inMemory bool) (*badger.DB, error) {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR)
	if inMemory {
		opts = opts.WithInMemory(true)
	}
	return badger.Open(opts)
}
