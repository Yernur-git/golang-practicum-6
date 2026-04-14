package handlers

import (
	"net/http"
	"practicum-6/config"
	"practicum-6/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	categoryFilter := strings.ToLower(c.Query("category"))
	authorIDFilter := c.Query("author_id")

	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	query := config.DB.Model(&models.Book{})

	if authorIDFilter != "" {
		query = query.Where("author_id = ?", authorIDFilter)
	}

	if categoryFilter != "" {
		var category models.Category
		if err := config.DB.Where("LOWER(name) = ?", categoryFilter).First(&category).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit, "total": 0, "data": []models.Book{}})
			return
		}
		query = query.Where("category_id = ?", category.ID)
	}

	var total int64
	query.Count(&total)

	var books []models.Book
	query.Offset((page - 1) * limit).Limit(limit).Find(&books)

	if books == nil {
		books = []models.Book{}
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  books,
	})
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func AddBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(book.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	if book.AuthorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_id must be positive"})
		return
	}
	if book.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id must be positive"})
		return
	}
	if book.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	var author models.Author
	if err := config.DB.First(&author, book.AuthorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author not found"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, book.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
		return
	}

	config.DB.Create(&book)
	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updated models.Book
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(updated.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	if updated.AuthorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_id must be positive"})
		return
	}
	if updated.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id must be positive"})
		return
	}
	if updated.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	config.DB.Model(&book).Updates(models.Book{
		Title:      updated.Title,
		AuthorID:   updated.AuthorID,
		CategoryID: updated.CategoryID,
		Price:      updated.Price,
	})

	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	config.DB.Delete(&book)
	c.Status(http.StatusNoContent)
}
