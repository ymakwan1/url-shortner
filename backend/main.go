package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to URL Shortener API",
		})
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
