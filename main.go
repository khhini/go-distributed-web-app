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

// HealthzHandler godoc
// describtion Health Check Handler
func HealthzHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "pong",
	})
}

func main() {
	router := gin.Default()
	router.GET("/", IndexHandler)
	router.GET("/healthz", IndexHandler)
	router.Run()
}
