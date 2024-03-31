package passport

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"

	"distrifs.dev/server/modules/globals"
)

type PassportKey struct {
	Key string `form:"key" binding:"required"`
}

func VerifyPassportRoute(c *gin.Context) {
	if !globals.Options_passport {
		c.JSON(403, gin.H{
			"success": false,
			"error":   "passport_not_enabled",
		})
		return
	}

	var passportKey PassportKey
	if err := c.ShouldBind(&passportKey); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_json_request",
		})
		return
	}

	err := globals.PermDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(passportKey.Key))

		if err != nil {
			c.JSON(404, gin.H{
				"success": false,
				"error":   "key_not_found",
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
			if err != nil {
				c.JSON(500, gin.H{
					"success": false,
					"error":   "database_valuecopy_error",
				})
				return nil
			}

			c.JSON(200, gin.H{
				"success": true,
				"data":    string(valCopy),
			})

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
}