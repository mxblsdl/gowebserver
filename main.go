package main

import (
	"fmt"
	"net/http"
	"webserver/internal/database"
	"webserver/internal/handlers"
	"webserver/internal/middleware"
)

func main() {
	// Initialize the database
	if err := database.InitDB(); err != nil {
		handlers.LogFatal("Failed to initialize the database: ", err)
	}

	mux := http.NewServeMux()

	// Login handlers
	mux.Handle("/", middleware.LoggingMiddleware(http.HandlerFunc(handlers.HomeHandler)))
	mux.Handle("/login", middleware.LoggingMiddleware(http.HandlerFunc(handlers.LoginHandler)))
	mux.Handle("/show_register", middleware.LoggingMiddleware(http.HandlerFunc(handlers.ShowRegisterPage)))
	mux.Handle("/register", middleware.LoggingMiddleware(http.HandlerFunc(handlers.RegisterHandler)))

	// apply recovery middleware
	handler := middleware.RecoveryMiddleware(mux)

	//serve static files
	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Configure the server
	port := ":8090"
	fmt.Printf("Server starting on http://localhost%s\n", port)

	// Start the server
	if err := http.ListenAndServe(port, handler); err != nil {
		handlers.LogFatal("Failed to server http: ", err)
	}
}
