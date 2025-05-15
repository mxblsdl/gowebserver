package middleware

import (
	// "log"
	"context"
	"net/http"
	"webserver/internal/database"
	"webserver/internal/logger"
)

// LoggingMiddleware logs the details of each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogInfo("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics and writes a 500 if there was one
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.LogError("Recovered from panic: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Key type for context values
type contextKey string

const UserDataKey contextKey = "userData"

func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogInfo("Retrieving API key from request header")
		// Log all headers
		// for key, value := range r.Header {
		// 	logger.LogDebug("Header: %s: %s", key, value)
		// }

		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			logger.LogWarning("No API key provided")
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		// Get user data associated with API key
		userData, err := database.GetUserByAPIKey(apiKey)
		if err != nil {
			logger.LogError("Error validating API key: %v", err)
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}
		logger.LogInfo("UserData: %v", userData)

		// Add user data to request context
		ctx := context.WithValue(r.Context(), UserDataKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
