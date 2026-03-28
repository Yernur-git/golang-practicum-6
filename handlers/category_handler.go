package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practicum-6/models"
	"strings"
)

var categories = make(map[int]models.Category)
var nextCategoryID = 1

func GetCategoryByID(id int) (models.Category, bool) {
	c, exists := categories[id]
	return c, exists
}

func GetCategories(c *gin.Context) {
	list := make([]models.Category, 0, len(categories))
	for _, cat := range categories {
		list = append(list, cat)
	}
	c.JSON(http.StatusOK, list)
}

func AddCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(category.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	category.ID = nextCategoryID
	nextCategoryID++
	categories[category.ID] = category

	c.JSON(http.StatusCreated, category)
}
