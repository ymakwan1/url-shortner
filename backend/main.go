package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ymakwan1/url-shortener/backend/handlers"
)

func main() {
	router := gin.Default()

	// Create short URL
	router.POST("/", handlers.CreateShortURL)

	// Get original URL
	router.GET("/:key", handlers.GetOriginalURL)

	router.Run(":3000")
}
