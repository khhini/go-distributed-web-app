package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// IndexHandler ...
func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "ping",
	})
}

// HealthzHandler godoc
// describtion Health Check Handler
func HealthzHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "ping",
	})
}

func main() {
	router := gin.Default()
	router.GET("/", IndexHandler)
	router.GET("/healthz", IndexHandler)
	router.Run()
}
