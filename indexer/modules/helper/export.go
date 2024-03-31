package helper

import (
	"fmt"
	"sort"

	"distrifs.dev/indexer/modules/globals"

	"github.com/dgraph-io/badger/v4"
)

// export badgerdb structure as a sorted JSON encoded structure

type AllKeysResponse struct {
	error bool
	data  map[string]string
}

func getAllKeys(db *badger.DB, res chan AllKeysResponse) {

	var data map[string]string

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("Processing key=%s, value=%s\n", k, v)
				data[string(k)] = string(v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	keyResponse := AllKeysResponse{}

	if err != nil {
		keyResponse.error = true
		res <- keyResponse
		return
	}

	keyResponse.error = false
	keyResponse.data = data
	res <- keyResponse
}

func Export() {
	db := globals.Badger_DB
	channel := make(chan AllKeysResponse, 1)

	go getAllKeys(db, channel)
	val := <-channel
	close(channel)

	// sort a-z

	data := val.data

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	sortedData := make(map[string]string, len(data))

	for _, k := range keys {
		sortedData[k] = data[k]
	}
}
