package handlers

import (
	"net/http"
	"practicum-6/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var books = make(map[int]models.Book)
var nextBookID = 1

func GetBooks(c *gin.Context) {
	categoryFilter := strings.ToLower(c.Query("category"))
	authorIDFilter := 0
	if raw := c.Query("author_id"); raw != "" {
		if id, err := strconv.Atoi(raw); err == nil {
			authorIDFilter = id
		}
	}

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

	var filtered []models.Book
	for _, book := range books {
		if authorIDFilter != 0 && book.AuthorID != authorIDFilter {
			continue
		}
		if categoryFilter != "" {
			cat, exists := GetCategoryByID(book.CategoryID)
			if !exists || strings.ToLower(cat.Name) != categoryFilter {
				continue
			}
		}
		filtered = append(filtered, book)
	}

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paginated := filtered[start:end]
	if paginated == nil {
		paginated = []models.Book{}
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  paginated,
	})
}

func GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, exists := books[id]
	if !exists {
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
	if book.AuthorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_id must be positive"})
		return
	}
	if book.CategoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id must be positive"})
		return
	}
	if book.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}
	if _, exists := GetAuthorByID(book.AuthorID); !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author not found"})
		return
	}
	if _, exists := GetCategoryByID(book.CategoryID); !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
		return
	}

	book.ID = nextBookID
	nextBookID++
	books[book.ID] = book

	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	if _, exists := books[id]; !exists {
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
	if updated.AuthorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_id must be positive"})
		return
	}
	if updated.CategoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id must be positive"})
		return
	}
	if updated.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	updated.ID = id
	books[id] = updated
	c.JSON(http.StatusOK, updated)
}

func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	if _, exists := books[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	delete(books, id)
	c.Status(http.StatusNoContent)
}
