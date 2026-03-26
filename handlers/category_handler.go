package handlers

import (
	"encoding/json"
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

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list := make([]models.Category, 0, len(categories))
	for _, c := range categories {
		list = append(list, c)
	}
	json.NewEncoder(w).Encode(list)
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(category.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	category.ID = nextCategoryID
	nextCategoryID++
	categories[category.ID] = category

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}
