package handlers

import (
	"encoding/json"
	"net/http"

	"flowtask/internal/models"
	"flowtask/internal/repository"
)

type CategoryHandler struct {
	repo *repository.SQLiteRepository
}

func NewCategoryHandler(repo *repository.SQLiteRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories, err := h.repo.GetCategories()
	if err != nil {
		http.Error(w, `{"error": "Failed to get categories"}`, http.StatusInternalServerError)
		return
	}
	if categories == nil {
		categories = []models.Category{}
	}
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, `{"error": "Invalid category name"}`, http.StatusBadRequest)
		return
	}

	cat, err := h.repo.CreateCategory(req.Name)
	if err != nil {
		http.Error(w, `{"error": "Category limit reached (max 10) or already exists"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func (h *CategoryHandler) HandleReorderCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid IDs"}`, http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateCategoriesOrder(req.IDs); err != nil {
		http.Error(w, `{"error": "Failed to update order"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order updated successfully"})
}
