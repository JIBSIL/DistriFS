package files

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"distrifs.dev/server/modules/helper"
)

type DirQuery struct {
	Dirname       string `form:"dirname" binding:"required"`
	PassportToken string `form:"token"`
}

type FileInfo struct {
	Name        string
	IsDirectory bool
}

func ReadDirRoute(c *gin.Context) {
	var dirQuery DirQuery
	if err := c.ShouldBind(&dirQuery); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	}

	var dirNameRaw = helper.SanitizePath(dirQuery.Dirname)
	var dirName string

	if dirNameRaw.Error {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	} else {
		dirName = dirNameRaw.Path
	}

	var dir, err = os.ReadDir(dirName)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "not_found",
		})
		return
	}

	files := make([]FileInfo, len(dir))
	files = files[:0]

	for _, fileRaw := range dir {
		var file = FileInfo{}

		var fileName = fileRaw.Name()

		var path = filepath.Join(dirName, fileName)
		if !helper.PermissionsCheck(c, dirQuery.PassportToken, path, false) {
			continue
		}

		file.Name = fileName
		file.IsDirectory = fileRaw.IsDir()

		files = append(files, file)
	}

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "json_error",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    files,
	})
}
