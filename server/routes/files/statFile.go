package files

import (
	"os"

	"github.com/gin-gonic/gin"

	"distrifs.dev/server/modules/helper"
)

func StatFileRoute(c *gin.Context) {
	// IF file exists: return stats about file
	// File doesn't exist: notify client

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

	var checksum string

	if statOutput.IsDir() {
		checksum = ""
	} else {
		checksum, err = helper.ChecksumFile(filename)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "checksum_error",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"success":      true,
		"name":         statOutput.Name(),
		"size":         statOutput.Size(),
		"isDirectory":  statOutput.IsDir(),
		"lastModified": statOutput.ModTime(),
		"hash":         checksum,
	})

}
