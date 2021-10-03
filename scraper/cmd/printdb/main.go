package main

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
)

func main() {

	db, err := badger.Open(badger.DefaultOptions("./databasedata"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			item.Value(func(v []byte) error {
				fmt.Printf("%v: %v\n", string(key), string(v))
				return nil
			})
		}

		return nil
	})
}
