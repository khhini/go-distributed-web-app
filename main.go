package main

import (
	"github.com/gin-gonic/gin"
)

// IndexHandler ...
func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "pong",
	})
}

func main() {
	router := gin.Default()
	router.GET("/", IndexHandler)
	router.Run()
}
