package helper

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"

	"distrifs.dev/server/modules/globals"
)

type ReturnPath struct {
	Error bool
	Path  string
}

func SanitizePath(path string) ReturnPath {

	var returned = ReturnPath{}

	// first, remove ..
	if strings.Contains(path, "..") {
		returned.Error = true
		returned.Path = ""
		return returned
	}

	// var fullpath string
	// if InternalOptions_pathSanitization {
	// 	fullpath = fmt.Sprintf("%s/%s", globals.Options_directory, path)
	// } else {
	// var fullpath = fmt.Sprintf("./%s", path)
	var fullpath = fmt.Sprintf("%s/%s", globals.Options_directory, path)
	// }

	returned.Error = false
	returned.Path = fullpath
	return returned
}

func VerifySig(data, signature []byte, pubkey *ecdsa.PublicKey) bool {
	//https://github.com/gtank/cryptopasta/blob/master/sign.go
	// hash message
	digest := sha256.Sum256(data)

	curveOrderByteSize := pubkey.Curve.Params().P.BitLen() / 8

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])

	return ecdsa.Verify(pubkey, digest[:], r, s)
}

func GetKeyByValue(db *badger.DB, valueRaw string) (string, error) {
	var key []byte
	value := []byte(valueRaw)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			itemValue, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			if bytes.Equal(itemValue, value) {
				key = item.KeyCopy(nil)
				return nil
			}
		}

		return badger.ErrKeyNotFound
	})

	return string(key), err
}

func PathEquality(p1 string, p2 string) bool {
	fi1, err1 := os.Lstat(p1)
	fi2, err2 := os.Lstat(p2)

	if err1 != nil || err2 != nil {
		fmt.Println("Error reading file info:", err1, err2)
		// assume paths are equal to throw an error in the passport logic
		return true
	}

	// Check if the files/directories are the same
	if os.SameFile(fi1, fi2) {
		return true
	} else {
		return false
	}
}

func CheckPassportForFile(requestedFile string) bool {
	if globals.Viper.GetBool("AllFilesPassportLocked") {
		return true
	}

	passportfiles := globals.Viper.Get("PassportProtectedFiles")

	// check if the file is in the passport
	for file := range passportfiles.(map[string](interface{})) {
		if PathEquality(file, requestedFile) {
			return true
		}
	}
	return false
}

// https://stackoverflow.com/questions/70585852/return-default-value-for-generic-type
func Zero[T any]() T {
	return *new(T)
}

func IsZero[T comparable](v T) bool {
	return v == *new(T)
}
