package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ymakwan1/url-shortener/backend/validator"
)

type ShortURL struct {
	Key      string `json:"key"`
	LongURLs string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

var ShortURLs = make(map[string]string)

func CreateShortURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validator.IsValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	for key, longUrl := range ShortURLs {
		if longUrl == req.URL {
			c.JSON(http.StatusOK, ShortURL{
				Key:      key,
				LongURLs: longUrl,
				ShortURL: "http://localhost/" + key,
			})
			return
		}
	}

	key := generateKey(len(ShortURLs) + 1)
	ShortURLs[key] = req.URL

	c.JSON(http.StatusCreated, ShortURL{
		Key:      key,
		LongURLs: req.URL,
		ShortURL: "http://localhost/" + key,
	})

}

func GetOriginalURL(c *gin.Context) {
	key := c.Param("key")
	longURL, ok := ShortURLs[key]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(http.StatusFound, longURL)
}

func generateKey(i int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	key := make([]byte, i)
	for v := range key {
		key[v] = letters[v%len(letters)]
	}
	return string(key)
}
