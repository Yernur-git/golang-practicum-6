package handlers

import (
	"encoding/json"
	"net/http"
	"practicum-6/models"
	"strings"
)

var authors = make(map[int]models.Author)
var nextAuthorID = 1

func GetAuthorByID(id int) (models.Author, bool) {
	a, exists := authors[id]
	return a, exists
}

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list := make([]models.Author, 0, len(authors))
	for _, a := range authors {
		list = append(list, a)
	}
	json.NewEncoder(w).Encode(list)
}

func AddAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(author.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	author.ID = nextAuthorID
	nextAuthorID++
	authors[author.ID] = author

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)
}
