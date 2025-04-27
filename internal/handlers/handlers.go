package handlers

import (
	"fmt"
	"net/http"
	"time"

	"webserver/internal/database"
	"webserver/internal/logger"
	"webserver/internal/models"
	"webserver/templates/pages"

	"golang.org/x/crypto/bcrypt"
)

// HomeHandler handles requests for the home page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := pages.Index(time.Now()).Render(r.Context(), w)
	if err != nil {
		logger.LogError("Error parsing template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO work on responses for user not found and incorrect password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	login := models.LoginRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	logger.LogInfo("Login request received for user: %s", login.Username)

	user, err := database.GetUser(login.Username)
	if err != nil {
		logger.LogWarning("User not found: %s", login.Username)
		w.Header().Set("HX-Trigger", fmt.Sprintf(`{"login" : {"username" : "%s", "type" : "error"}}`, login.Username))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		logger.LogWarning("Incorrect password for user: %s", login.Username)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("HX-Trigger", fmt.Sprintf(`{"login" : {"username" : "%s", "type" : "error"}}`, login.Username))
		// w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	data := pages.PageData{
		Username: login.Username,
		Key:      user.APIKey,
		FolderId: user.FolderId,
	}

	logger.LogInfo("Logged in user: %s", login.Username)
	err = pages.Main(data).Render(r.Context(), w)

	if err != nil {
		logger.LogError("Error parsing template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogInfo("User logged in successfully: %s", login.Username)

}

func ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := pages.Register().Render(r.Context(), w)
	fmt.Println("Show register page request received!!", err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Register request received!!")
	err := r.ParseForm()
	if err != nil {
		logger.LogError("Error parsing form: ", err)
		return
	}

	register := models.LoginRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	exists := database.UsernameExists(register.Username)
	if exists {
		logger.LogWarning("Username already exists: %s", register.Username)
		w.Header().Set("HX-Trigger", fmt.Sprintf(`{"register" : {"username" : "%s", "type" : "error"}}`, register.Username))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	register.Password = string(hashedPassword)

	err = database.CreateUser(register.Username, register.Password)
	if err != nil {
		logger.LogError("Error creating user: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.LogInfo("User created successfully: %s", register.Username)
	w.Header().Set("HX-Trigger", fmt.Sprintf(`{"register" : {"username" : "%s", "type" : "success"}}`, register.Username))
	w.Write([]byte(""))
}

func FilePathHandler(w http.ResponseWriter, r *http.Request) {
	// TODO implement file path handler 
	logger.LogInfo("File path request received")
}


func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO implement items handler
	logger.LogInfo("Items request received")
}
