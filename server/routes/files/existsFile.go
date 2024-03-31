package files

import (
	"os"

	"github.com/gin-gonic/gin"

	"distrifs.dev/server/modules/helper"
)

func ExistsRoute(c *gin.Context) {
	// simple bool for lightweight clients
	// return true if file exists, else return false
	var fileQuery FileQuery
	if err := c.ShouldBind(&fileQuery); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	}

	var filenameRaw = helper.SanitizePath(fileQuery.Filename)
	var filename string

	if filenameRaw.Error {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	} else {
		filename = filenameRaw.Path
	}

	if !helper.PermissionsCheck(c, fileQuery.PassportToken, filename, true) {
		return
	}

	var statOutput, err = os.Stat(filename)

	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "not_found",
		})
		return
	}

	if statOutput.IsDir() {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "is_directory",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}
