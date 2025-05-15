package main

import (
	"flag"
	"fmt"
	"net/http"
	"webserver/internal/database"
	"webserver/internal/handlers"
	"webserver/internal/logger"
	"webserver/internal/middleware"
	"webserver/pkg/config"
)

func main() {
	devMode := flag.Bool("dev", false, "Run in development mode")
	flag.Parse()

	config.DevMode = *devMode

	// Initialize the logger
	if err := logger.InitLogger("DEBUG"); err != nil {
		logger.LogFatal("Failed to initialize logger: ", err)
	}

	// Initialize the database
	if err := database.InitDB(); err != nil {
		logger.LogFatal("Failed to initialize the database: ", err)
	}

	mux := http.NewServeMux()

	protected := func(handler http.HandlerFunc) http.Handler {
		return middleware.LoggingMiddleware(middleware.RequireAPIKey(http.HandlerFunc(handler)))
	}

	// Login handlers
	if config.DevMode {
		mux.HandleFunc("/", handlers.DevHomeHandler)
	} else {
		mux.Handle("/", middleware.LoggingMiddleware(http.HandlerFunc(handlers.HomeHandler)))
	}

	mux.Handle("/login", middleware.LoggingMiddleware(http.HandlerFunc(handlers.LoginHandler)))
	mux.Handle("/show_register", middleware.LoggingMiddleware(http.HandlerFunc(handlers.ShowRegisterPage)))
	mux.Handle("/register", middleware.LoggingMiddleware(http.HandlerFunc(handlers.RegisterHandler)))

	// API handlers
	mux.Handle("/filepath", protected(handlers.FilePathHandler))
	mux.Handle("/items", protected(handlers.ItemsHandler))
	mux.Handle("/upload", protected(handlers.UploadHandler))
	mux.Handle("/download", protected(handlers.DownloadHandler))

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
		logger.LogFatal("Failed to server http: ", err)
	}
}
