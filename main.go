package main

import (
	"log"
	"net/http"

	"flowtask/internal/database"
	"flowtask/internal/handlers"
	"flowtask/internal/repository"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db := database.InitDB("flowtask.db")
	defer db.Close()

	repo := repository.NewSQLiteRepository(db)
	taskHandle := handlers.NewTaskHandler(repo)
	categoryHandler := handlers.NewCategoryHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/tasks", taskHandle.HandleGetTasks)
	mux.HandleFunc("POST /api/tasks", taskHandle.HandleCreateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", taskHandle.HandleDeleteTask)
	mux.HandleFunc("PATCH /api/tasks/{id}/toggle", taskHandle.HandleToggleTask)
	mux.HandleFunc("GET /api/categories", categoryHandler.HandleGetCategories)
	mux.HandleFunc("POST /api/categories", categoryHandler.HandleCreateCategory)
	mux.HandleFunc("PATCH /api/categories/reorder", categoryHandler.HandleReorderCategories)

	handlerWithCORS := enableCORS(mux)

	log.Println("Server running http://localhost:8080")
	if err := http.ListenAndServe(":8080", handlerWithCORS); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
