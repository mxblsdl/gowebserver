package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"webserver/internal/database"
	"webserver/internal/models"
	"webserver/templates/pages"

	"golang.org/x/crypto/bcrypt"
)

// HomeHandler handles requests for the home page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := pages.Index(time.Now()).Render(r.Context(), w)
	if err != nil {
		LogError("Error parsing template: ", err)
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

	user, err := database.GetUser(login.Username)
	if err != nil {
		w.Header().Set("HX-Trigger", `{"login" : "User not found"}`)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)

		return
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("HX-Trigger", fmt.Sprintf(`{"login" : {"username" : "%s", "type" : "error"}}`, login.Username))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := pages.PageData{
		Username: login.Username,
		Key:      "placeholder",
		FolderId: "placeholder",
	}
	err = pages.Main(data).Render(r.Context(), w)
	if err != nil {
		LogError("Error parsing template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("User logged in successfully")

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
		LogError("Error parsing form: ", err)
		return
	}

	register := models.LoginRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	exists := database.UsernameExists(register.Username)
	if exists {
		log.Println("Username already exists: ", register.Username)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("User created successfully")
	w.Header().Set("HX-Trigger", fmt.Sprintf(`{"register" : {"username" : "%s", "type" : "success"}}`, register.Username))
	w.Write([]byte(""))
}
