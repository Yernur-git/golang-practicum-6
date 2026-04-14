package handlers

import (
	"net/http"
	"practicum-6/config"
	"practicum-6/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAuthors(c *gin.Context) {
	var authors []models.Author
	config.DB.Find(&authors)
	c.JSON(http.StatusOK, authors)
}

func AddAuthor(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(author.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	config.DB.Create(&author)
	c.JSON(http.StatusCreated, author)
}
