package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/TomaszSkrzyp/go-lang-api/toDo/internal/auth"
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

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizarionHeader := r.Header.Get("Authorization")
		if authorizarionHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing Authorization header"})
			return
		}

		parts := strings.Split(authorizarionHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Authorization header format"})
			return
		}
		tokenStr := parts[1]
		_, err := auth.ValidateToken(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
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
	// Protected routes - require valid JWT token
	router.Handle("/api/todos/{id}", jwtMiddleware(http.HandlerFunc(storage.HandleGet))).Methods("GET")
	router.Handle("/api/todos/{id}", jwtMiddleware(http.HandlerFunc(storage.HandleUpdateTask))).Methods("PATCH")
	router.Handle("/api/todos/{id}", jwtMiddleware(http.HandlerFunc(storage.HandleRemove))).Methods("DELETE")
	router.Handle("/api/todos", jwtMiddleware(http.HandlerFunc(storage.HandleAdd))).Methods("POST")
	router.Handle("/api/todos", jwtMiddleware(http.HandlerFunc(storage.HandleGetAll))).Methods("GET")

	router.HandleFunc("/login", dbControl.LoginHandler).Methods("POST")
	// Serve static frontend (React build output, for example)
	fs := http.FileServer(http.Dir("./todo-frontend/build"))
	router.PathPrefix("/").Handler(fs)

	// Wrap router with CORS middleware and start the server
	handlersWithCors := enableCors(router)
	http.ListenAndServe(":8090", handlersWithCors)
}
