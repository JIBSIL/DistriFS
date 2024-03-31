package routes

import (
	"distrifs.dev/indexer/modules/globals"
	"github.com/gin-gonic/gin"
)

type ExplicitQuery struct {
	Hash   string `form:"hash" binding:"required"`
	Server string `form:"server" binding:"required"`
}

func getFileWalk(v map[string]globals.HashItem, hash string) globals.HashItem {
	// fmt.Printf("Visiting %v\n", v)
	returned := globals.HashItem{}
out:
	for _, v := range v {
		if v.Hash == hash {
			returned = v
		}
		if v.SubFiles != nil {
			val := getFileWalk(v.SubFiles, hash)
			if val.Hash != "" {
				returned = val
				break out
			}
		}
	}
	return returned
}

func GetFileRoute(c *gin.Context) {
	// get file by hash

	var hashQuery ExplicitQuery
	if err := c.ShouldBind(&hashQuery); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	}

	server := globals.IndexerData[hashQuery.Server]

	if !server.Online {
		c.JSON(500, gin.H{
			"success": false,
			"data":    "server_not_online",
		})
		return
	}

	file := getFileWalk(server.Files.SubFiles, hashQuery.Hash)

	if file.Hash != hashQuery.Hash {
		c.JSON(500, gin.H{
			"success": false,
			"data":    "no_file_found",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    file,
	})
}
