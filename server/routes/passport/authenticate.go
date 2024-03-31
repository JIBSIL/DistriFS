package passport

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"distrifs.dev/server/modules/globals"
)

// todo: rate limiting here

func AuthenticateRoute(c *gin.Context) {

	if !globals.Options_passport {
		c.JSON(403, gin.H{
			"success": false,
			"error":   "passport_not_enabled",
		})
		return
	}

	// give user a random uuid to sign
	key := uuid.New().String()
	value := uuid.New().String()

	// new writable transaction
	txn := globals.PassportDB.NewTransaction(true)
	defer txn.Discard()

	err := txn.Set([]byte(key), []byte(value))
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "error_setting_in_db",
		})
		return
	}

	// reference:
	// KEY is the user's identifier for this transaction
	// VALUE is what should be signed

	// Commit the transaction and check for error.
	if err := txn.Commit(); err != nil {
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "error_commiting_db",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"key":   key,
			"value": value,
		},
	})
}