package controller

import "github.com/gin-gonic/gin"

func Ping(c *gin.Context) {
	successResponse(c, 200, 1, gin.H{
		"message": "pong",
	})
}
