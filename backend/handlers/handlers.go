package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShortURL struct {
	Key      string `json:"key"`
	LongURLs string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

var ShortURLs = make(map[string]string)

func CreateShortURL(c *gin.Context) {
	var req struct {
		URL string `json: "url" binding: "required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func generateKey(i int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY0123456789"
	key := make([]byte, i)
	for v := range key {
		key[v] = letters[v%len(letters)]
	}
	return string(key)
}
