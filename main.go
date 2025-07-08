package main

import (
	"fmt"
	"log"
	"net/http"
	"todoApp/auth"
	"todoApp/database"
	"todoApp/handlers"

	"github.com/gorilla/mux"
)

func main() {
	//initialise database
	database.InitDB()
	//routes
	r := mux.NewRouter()

	//public routes
	r.HandleFunc("/register", handlers.RegisterUser)
	r.HandleFunc("/login", handlers.Login)

	//sub router for protected routes
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(auth.AuthMiddleware)

	//authenticated routes
	protected.HandleFunc("/protect", handlers.ProtectedHandler).Methods("GET")
	protected.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")
	protected.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
	protected.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")
	protected.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE")

	//start server
	fmt.Println("starting server on http/localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("error connecting to server: ", err)
	}
}
