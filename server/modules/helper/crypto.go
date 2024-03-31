package helper

import (
	"crypto"
	"encoding/hex"
	"io"
	"os"

	"distrifs.dev/server/modules/globals"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
)

// this whole file is passport routes stuff
// with all the gin-specific code, it didn't fit into helper.go

func GetKey(c *gin.Context, key string) string {
	var value string
	err := globals.PermDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))

		if err != nil {
			c.JSON(404, gin.H{
				"success": false,
				"error":   "key_not_found",
			})
			return err
		}

		err = item.Value(func(val []byte) error {
			return err
		})

		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "database_internal_error",
			})
			return err
		} else {
			var valCopy []byte
			valCopy, err = item.ValueCopy(nil)
			if err != nil {
				c.JSON(500, gin.H{
					"success": false,
					"error":   "database_valuecopy_error",
				})
				return err
			}

			value = string(valCopy)

			return nil
		}
	})

	if err != nil {
		// removed for now as this only is necessary in specific cases
		// c.JSON(500, gin.H{
		// 	"success": false,
		// 	"error":   "database_internal_error",
		// })
		return Zero[string]()
	}
	return value
}

func VerifyPassport(c *gin.Context, key string, requestedFile string) bool {
	// get user's creds
	var userPubkey = GetKey(c, key)
	if IsZero[string](userPubkey) {
		return false
	}

	// test if user is admin
	var adminkeys = globals.Viper.GetStringSlice("AdminKeys")
	for _, adminkey := range adminkeys {
		if userPubkey == adminkey {
			return true
		}
	}

	// test if user has perms on file

	// but we need to get the file properly first
	passportfiles := globals.Viper.GetStringMapStringSlice("PassportProtectedFiles")
	keys := []string{}

	for file, fileKeys := range passportfiles {
		if PathEquality(file, requestedFile) {
			keys = fileKeys
			break
		}
	}

	for _, fileKey := range keys {
		if userPubkey == fileKey {
			return true
		}
	}

	return false
}

// three line permissions check
func PermissionsCheck(c *gin.Context, key string, requestedFile string, throwErrors bool) bool {
	inPassport := CheckPassportForFile(requestedFile)

	if inPassport {
		if key != "" {
			isVerified := VerifyPassport(c, key, requestedFile)
			if !isVerified {
				// we can assume that an error was already thrown in c, and hence, we don't want to throw another
				return false
			} else {
				return true
			}
		} else {
			if throwErrors {
				c.JSON(403, gin.H{
					"success": false,
					"error":   "passport_required",
				})
			}
			return false
		}
	} else {
		return true
	}
}

func ChecksumFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := crypto.SHA256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
