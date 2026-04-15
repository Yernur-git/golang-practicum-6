package main

import (
	"log"
	"practicum-6/config"
	"practicum-6/handlers"
	"practicum-6/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.GET("/authors", handlers.GetAuthors)
	r.POST("/authors", handlers.AddAuthor)

	r.GET("/books", handlers.GetBooks)
	r.POST("/books", handlers.AddBook)
	r.GET("/books/:id", handlers.GetBookByID)
	r.PUT("/books/:id", handlers.UpdateBook)
	r.DELETE("/books/:id", handlers.DeleteBook)

	r.GET("/categories", handlers.GetCategories)
	r.POST("/categories", handlers.AddCategory)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/books/favorites", handlers.GetFavorites)
		auth.PUT("/books/:id/favorites", handlers.AddToFavorites)
		auth.DELETE("/books/:id/favorites", handlers.RemoveFromFavorites)
	}

	log.Println("Server running on :8080")
	r.Run(":8080")
}
