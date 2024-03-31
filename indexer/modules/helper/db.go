package helper

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"distrifs.dev/indexer/modules/globals"
)

type DbReadReturn struct {
	Error bool
	Data  string
}

func DbReadValue(key string) DbReadReturn {
	var answer []byte

	err := globals.Badger_DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))

		if err != nil {
			return err
		}

		var valCopy []byte
		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})

		if err != nil {
			return err
		}

		fmt.Printf("The answer is: %s\n", valCopy)

		_, err = item.ValueCopy(answer)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		returned := DbReadReturn{}
		returned.Error = true
		return returned
	}

	returned := DbReadReturn{}
	returned.Error = false
	returned.Data = string(answer)
	return returned
}
