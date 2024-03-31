package routes

import (
	"distrifs.dev/indexer/modules/globals"
	"github.com/gin-gonic/gin"
)

func GetAllFilesRoute(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    globals.IndexerData,
	})
}
