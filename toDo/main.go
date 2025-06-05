package main

import (
	"net/http"

	"github.com/TomaszSkrzyp/go-lang-api/toDo/internal/dbControl"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
	db := dbControl.InitDB()
	defer db.Close()

	storage := dbControl.TodoStorage{
		DB: db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/todos/{id}", storage.HandleGet).Methods("GET")
	router.HandleFunc("/todos/{id}", storage.HandleUpdateTask).Methods("PATCH")
	router.HandleFunc("/todos/{id}", storage.HandleRemove).Methods("DELETE")
	router.HandleFunc("/todos", storage.HandleAdd).Methods("POST")
	router.HandleFunc("/todos", storage.HandleGetAll).Methods("GET")

	fs := http.FileServer(http.Dir("./todo-frontend/build"))
	router.PathPrefix("/").Handler(fs)
	
	handlersWithCors := enableCors(router)
	http.ListenAndServe(":8090", handlersWithCors)
}
