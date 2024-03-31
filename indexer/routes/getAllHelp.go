package routes

import (
	"github.com/gin-gonic/gin"
)

func GetAllHelpRoute(c *gin.Context) {
	// getAll doesn't work without more args
	c.JSON(400, gin.H{
		"success": false,
		"error":   "specify_file_to_download",
	})
}
