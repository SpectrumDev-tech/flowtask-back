// Package handlers
package handlers

import (
	"encoding/json"
	"flowtask/internal/models"
	"flowtask/internal/repository"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	repo *repository.SQLiteRepository
}

func NewTaskHandler(repo *repository.SQLiteRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	energyFilter := r.URL.Query().Get("energy")

	tasks, err := h.repo.GetAll(energyFilter)
	if err != nil {
		http.Error(w, `{"error": "Failed to get tasks"}`, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Title       string             `json:"title"`
		Description string             `json:"description"`
		Energy      models.EnergyLevel `json:"energy"`
		Category    string             `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, `{"error": "Empty title"}`, http.StatusBadRequest)
		return
	}

	if req.Energy == "" {
		req.Energy = models.EnergyMedium
	}
	if req.Category == "" {
		req.Category = "code"
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Energy:      req.Energy,
		Category:    req.Category,
	}
	if err := h.repo.Create(task); err != nil {
		http.Error(w, `{"error": "Error saving task"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID of task"}`, http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, `{"error": "Delete task failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "the task was successfully deleted"})
}

func (h *TaskHandler) HandleToggleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID of task"}`, http.StatusBadRequest)
		return
	}

	if err := h.repo.ToggleComplete(id); err != nil {
		http.Error(w, `{"error": "The task status could not be updated"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "The task status has been successfully updated"})
}
