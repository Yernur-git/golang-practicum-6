package handlers

import (
	"net/http"
	"practicum-6/config"
	"practicum-6/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFavorites(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

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

	var total int64
	config.DB.Model(&models.FavoriteBook{}).Where("user_id = ?", userID).Count(&total)

	var favorites []models.FavoriteBook
	config.DB.Where("user_id = ?", userID).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&favorites)

	bookIDs := make([]uint, 0, len(favorites))
	for _, f := range favorites {
		bookIDs = append(bookIDs, f.BookID)
	}

	var books []models.Book
	if len(bookIDs) > 0 {
		config.DB.Where("id IN ?", bookIDs).Find(&books)
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  books,
	})
}

func AddToFavorites(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := config.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var existing models.FavoriteBook
	result := config.DB.Where("user_id = ? AND book_id = ?", userID, bookID).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Book is already in favorites"})
		return
	}

	favorite := models.FavoriteBook{
		UserID: userID,
		BookID: uint(bookID),
	}
	config.DB.Create(&favorite)

	c.JSON(http.StatusOK, gin.H{"message": "Book added to favorites"})
}

func RemoveFromFavorites(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	result := config.DB.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.FavoriteBook{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
