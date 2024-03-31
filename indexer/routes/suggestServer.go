package routes

import (
	"fmt"

	"distrifs.dev/indexer/modules/globals"

	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PutFileQuery struct {
	Server string `form:"server" binding:"required"`
}

func PutFileRoute(c *gin.Context) {
	// suggest a new file to the indexer
	// THIS SHOULD BE RATE-LIMITED IN PRODUCTION!!!

	var query PutFileQuery
	if err := c.ShouldBind(&query); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	}

	// first make sure the server is running
	res, err := http.Get(query.Server)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"data":    "http_error",
		})
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"data":    "error_reading_http_body",
		})
		return
	}

	var parsed gin.H
	json.Unmarshal(body, &parsed)

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"data":    "json_parsing_error",
		})
		return
	}

	if parsed["alive"] == true && parsed["server"] == "DistriFS Distributed File Server" {
		// check for living server in the system
		if !globals.IndexerData[query.Server].Online {
			fmt.Println("Adding server...")
			server := globals.IndexerServer{}

			server.Online = true
			server.Hashes = 0
			server.Files = globals.HashItem{}
		}
	}
}
