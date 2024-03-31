package passport

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"distrifs.dev/server/modules/globals"
	"distrifs.dev/server/modules/helper"
)

type SignedMessageKey struct {
	Pubkey        string `form:"pubkey" binding:"required"`
	SignedMessage string `form:"signedmessage" binding:"required"`
	MessageKey    string `form:"messagekey" binding:"required"`
}

func GetKeyRoute(c *gin.Context) {
	if !globals.Options_passport {
		c.JSON(403, gin.H{
			"success": false,
			"error":   "passport_not_enabled",
		})
		return
	}

	var messageKey SignedMessageKey
	if err := c.ShouldBind(&messageKey); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_json_request",
		})
		return
	}

	var value []byte

	err := globals.PassportDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(messageKey.MessageKey))

		if err != nil {
			c.JSON(404, gin.H{
				"success": false,
				"error":   "no_key_found",
			})
			return nil
		}

		var valCopy []byte
		err = item.Value(func(val []byte) error {
			return nil
		})

		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "database_internal_error",
			})
			return nil
		} else {
			valCopy, err = item.ValueCopy(nil)

			value = valCopy

			if err != nil {
				c.JSON(500, gin.H{
					"success": false,
					"error":   "database_valuecopy_error",
				})
				return nil
			}

			return nil
		}
	})

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "database_internal_error",
		})
		return
	}

	pubkeyBytes, err := hex.DecodeString(messageKey.Pubkey)

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "pubkey_decoding_err",
		})
		return
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), pubkeyBytes)
	publicKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	signedMessage, err := hex.DecodeString(messageKey.SignedMessage)

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "signedmessage_decoding_err",
		})
		return
	}

	sig := helper.VerifySig([]byte(value), signedMessage, publicKey)
	if sig {
		keyRaw := uuid.New().String()
		key := strings.ReplaceAll(keyRaw, "-", "")

		txn := globals.PermDB.NewTransaction(true)
		defer txn.Discard()

		// check if a key already exists
		_, err := helper.GetKeyByValue(globals.PermDB, string(messageKey.Pubkey))
		if err == nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "key_already_exists",
			})
			return
		}

		err = txn.Set([]byte(key), []byte(messageKey.Pubkey))
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "error_setting_in_db",
			})
			return
		}

		if err := txn.Commit(); err != nil {
			if err != nil {
				c.JSON(500, gin.H{
					"success": false,
					"error":   "error_commiting_db",
				})
				return
			}
		}

		txn = globals.PassportDB.NewTransaction(true)
		err = txn.Delete([]byte(messageKey.MessageKey))
		if err != nil {
			return
		}

		if err := txn.Commit(); err != nil {
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    key,
		})
		return
	}
}
