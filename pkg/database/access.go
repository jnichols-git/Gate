package database

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

// Database "generics"; get/set from collection of badger databases

type DatabaseCollection map[string]*badger.DB

// Open a database at /tmp/{name}. if write == false, NO DATA will be written to disk.
// Leave write = true unless you're testing.
func openDatabase(name string, write bool) (*badger.DB, error) {
	if write {
		return badger.Open(badger.DefaultOptions("./dat/database/" + name))
	} else {
		return badger.Open(badger.DefaultOptions("").WithInMemory(true))
	}
}

// Open every DB in a DatabaseCollection. See description for write arg above in openDatabase.
func (coll DatabaseCollection) opened(write bool) (DatabaseCollection, error) {
	var err error
	for name := range coll {
		coll[name], err = openDatabase(name, write)
		if err != nil {
			return nil, err
		}
	}
	return coll, nil
}

func setKey(db *badger.DB, key, value []byte) error {
	if db == nil {
		return errors.New("Attempted to setKey on a nil database")
	}
	if err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err == badger.ErrConflict {
			return err
		}
		if err == badger.ErrTxnTooBig {
			return err
		}
		return nil
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func getKey(db *badger.DB, key []byte) ([]byte, error) {
	if db == nil {
		return nil, errors.New("Attempted to getKey on a nil database")
	}
	var val []byte
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		item.Value(func(v []byte) error {
			val = v
			return nil
		})
		return nil
	}); err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func keyExists(db *badger.DB, key []byte) (bool, error) {
	if db == nil {
		return false, errors.New("Attempted to getKey on a nil database")
	}
	var val bool = false
	if err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			if string(it.Item().Key()) == string(key) {
				val = true
			}
		}
		return nil
	}); err != nil {
		return false, err
	} else {
		return val, nil
	}
}
