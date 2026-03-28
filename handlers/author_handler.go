package handlers

import (
	"net/http"
	"practicum-6/models"
	"strings"

	"github.com/gin-gonic/gin"
)

var authors = make(map[int]models.Author)
var nextAuthorID = 1

func GetAuthorByID(id int) (models.Author, bool) {
	a, exists := authors[id]
	return a, exists
}

func GetAuthors(c *gin.Context) {
	list := make([]models.Author, 0, len(authors))
	for _, a := range authors {
		list = append(list, a)
	}
	c.JSON(http.StatusOK, list)
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

	author.ID = nextAuthorID
	nextAuthorID++
	authors[author.ID] = author

	c.JSON(http.StatusCreated, author)
}
