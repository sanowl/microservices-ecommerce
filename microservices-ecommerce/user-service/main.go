package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = make(map[string]User)

func init() {
	// Initializing a few users for demonstration
	users["1"] = User{ID: "1", Name: "John Doe", Email: "john@example.com"}
	users["2"] = User{ID: "2", Name: "Jane Doe", Email: "jane@example.com"}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, ok := users[params["id"]]
	if ok {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users[user.ID] = user
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := users[params["id"]]; !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	users[params["id"]] = updatedUser
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := users[params["id"]]; !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(users, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Set up logging
	log.SetOutput(os.Stdout)

	// Create a new router
	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Get the port from the environment variables, default to 8081 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start the server
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
