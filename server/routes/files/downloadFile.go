package files

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"

	"distrifs.dev/server/modules/globals"
)

type DownloadQuery struct {
	ID string `form:"id" binding:"required"`
}

func DownloadFileRoute(c *gin.Context) {

	// parse uri
	var downloadQuery DownloadQuery
	if err := c.ShouldBind(&downloadQuery); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	}

	// host a new one-time route with a key
	err := globals.Badger_DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(downloadQuery.ID))

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

			filepath := fmt.Sprintf("./%s", string(valCopy))
			individualFilepathSlice := strings.Split(string(valCopy), "/")
			individualFilepath := individualFilepathSlice[len(individualFilepathSlice)-1]

			byteFile, err := os.ReadFile(filepath)
			if err != nil {
				fmt.Println(err)
			}

			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", individualFilepath))
			c.Data(200, "application/octet-stream", byteFile)

			if globals.Options_oneTimeFiles {
				// delete file from db

				// need to setup a new read-write tx for this
				txn = globals.Badger_DB.NewTransaction(true)
				err := txn.Delete([]byte(downloadQuery.ID))
				if err != nil {
					return nil
				}

				if err := txn.Commit(); err != nil {
					return nil
				}
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
}