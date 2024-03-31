package routes

import (
	"distrifs.dev/indexer/modules/globals"
	"github.com/gin-gonic/gin"
)

func GetAllFormatted(c *gin.Context) {
	// get all known indexers on the network (for discovery)

	var indexerList []string

	for k, v := range globals.IndexerData {
		if v.Online {
			indexerList = append(indexerList, k)
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    indexerList,
	})
}
