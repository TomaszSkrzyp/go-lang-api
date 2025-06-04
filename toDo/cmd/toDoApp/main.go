package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db := initDB()
	defer db.Close()

	storage := todo_storage{
		db: db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/todos/{id}", storage.handleGet).Methods("GET")
	router.HandleFunc("/todos/{id}", storage.handleUpdateTask).Methods("PUT")
	router.HandleFunc("/todos/{id}", storage.handleRemove).Methods("DELETE")
	router.HandleFunc("/todos", storage.handleAdd).Methods("POST")
	router.HandleFunc("/todos", storage.handleGetAll).Methods("GET")

	handlersWithCors := enableCors(router)
	http.ListenAndServe(":8090", handlersWithCors)
}
