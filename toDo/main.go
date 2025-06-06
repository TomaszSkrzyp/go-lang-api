package main

import (
	"net/http"

	"github.com/TomaszSkrzyp/go-lang-api/toDo/internal/dbControl"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// enableCors is middleware that sets CORS headers for local development,
// allowing frontend on localhost:3000 to access backend on :8090.
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
	// Initialize and defer-close the database connection
	db := dbControl.InitDB()
	defer db.Close()

	// Initialize data layer with DB connection
	storage := dbControl.TodoStorage{DB: db}

	// Set up the router and endpoint handlers
	router := mux.NewRouter()
	router.HandleFunc("/todos/{id}", storage.HandleGet).Methods("GET")
	router.HandleFunc("/todos/{id}", storage.HandleUpdateTask).Methods("PATCH")
	router.HandleFunc("/todos/{id}", storage.HandleRemove).Methods("DELETE")
	router.HandleFunc("/todos", storage.HandleAdd).Methods("POST")
	router.HandleFunc("/todos", storage.HandleGetAll).Methods("GET")

	// Serve static frontend (React build output, for example)
	fs := http.FileServer(http.Dir("./todo-frontend/build"))
	router.PathPrefix("/").Handler(fs)

	// Wrap router with CORS middleware and start the server
	handlersWithCors := enableCors(router)
	http.ListenAndServe(":8090", handlersWithCors)
}
