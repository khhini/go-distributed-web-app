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
