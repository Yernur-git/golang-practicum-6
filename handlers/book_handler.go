package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"practicum-6/models"
	"strconv"
	"strings"
)

var books = make(map[int]models.Book)
var nextBookID = 1

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()

	categoryFilter := strings.ToLower(q.Get("category"))
	authorIDFilter := 0
	if raw := q.Get("author_id"); raw != "" {
		if id, err := strconv.Atoi(raw); err == nil {
			authorIDFilter = id
		}
	}

	page := 1
	limit := 10
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := q.Get("limit"); l != "" {
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

	response := map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  paginated,
	}
	json.NewEncoder(w).Encode(response)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, exists := books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(book.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if book.AuthorID <= 0 {
		http.Error(w, "author_id is required and must be positive", http.StatusBadRequest)
		return
	}
	if book.CategoryID <= 0 {
		http.Error(w, "category_id is required and must be positive", http.StatusBadRequest)
		return
	}
	if book.Price <= 0 {
		http.Error(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}

	if _, exists := GetAuthorByID(book.AuthorID); !exists {
		http.Error(w, "author not found", http.StatusBadRequest)
		return
	}
	if _, exists := GetCategoryByID(book.CategoryID); !exists {
		http.Error(w, "category not found", http.StatusBadRequest)
		return
	}

	book.ID = nextBookID
	nextBookID++
	books[book.ID] = book

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if _, exists := books[id]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	var updated models.Book
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(updated.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if updated.AuthorID <= 0 {
		http.Error(w, "author_id is required and must be positive", http.StatusBadRequest)
		return
	}
	if updated.CategoryID <= 0 {
		http.Error(w, "category_id is required and must be positive", http.StatusBadRequest)
		return
	}
	if updated.Price <= 0 {
		http.Error(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}

	updated.ID = id
	books[id] = updated
	json.NewEncoder(w).Encode(updated)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if _, exists := books[id]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	delete(books, id)
	w.WriteHeader(http.StatusNoContent)
}
