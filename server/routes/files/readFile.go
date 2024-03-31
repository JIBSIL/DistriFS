package files

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"distrifs.dev/server/modules/globals"
	"distrifs.dev/server/modules/helper"
)

type FileQuery struct {
	Filename      string `form:"filename" binding:"required"`
	PassportToken string `form:"token"`
}

func ReadFileRoute(c *gin.Context) {
	// return hash
	var fileQuery FileQuery
	if err := c.ShouldBind(&fileQuery); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	}

	var staticFileRaw = helper.SanitizePath(fileQuery.Filename)
	var staticFile string

	if staticFileRaw.Error {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	} else {
		staticFile = staticFileRaw.Path
	}

	if !helper.PermissionsCheck(c, fileQuery.PassportToken, staticFile, true) {
		return
	}

	var file, err = os.Open(staticFile)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "not_found",
		})
		return
	}

	fileStats, err := file.Stat()

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "stat_error",
		})
		return
	}

	if fileStats.IsDir() {
		if err != nil {
			c.JSON(404, gin.H{
				"success": false,
				"error":   "is_directory",
			})
			return
		}
	}

	// safe to provide a download now

	// Start a writable transaction.
	txn := globals.Badger_DB.NewTransaction(true)
	defer txn.Discard()

	// create uuid for the user to use
	uuid := uuid.New().String()

	err = txn.Set([]byte(uuid), []byte(staticFile))
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "error_setting_in_db",
		})
		return
	}

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
		"data":    uuid,
	})
}
