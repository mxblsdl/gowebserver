package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"webserver/internal/database"
	"webserver/internal/logger"
	"webserver/internal/middleware"
	"webserver/internal/models"
	"webserver/pkg/config"
	"webserver/templates/components"
	"webserver/templates/pages"

	"strconv"

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

func DevHomeHandler(w http.ResponseWriter, r *http.Request) {
	data := pages.PageData{
		Username: config.DevUser.Username,
		Key:      config.DevUser.APIKey,
		FolderId: config.DevUser.FolderId,
	}

	err := pages.Main(data).Render(r.Context(), w)
	if err != nil {
		logger.LogError("Error rendering dev page: %v", err)
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
	logger.LogInfo("File path request received")
	userData, ok := r.Context().Value(middleware.UserDataKey).(database.UserData)
	logger.LogDebug("User data retrieved from context: %v", userData)
	if !ok || userData.UserId == 0 {
		logger.LogError("User data not found or invalid in context")
		http.Error(w, "User data not found", http.StatusInternalServerError)
		return
	}

	folderId, err := extractFolderId(r)
	if err != nil {
		logger.LogError("Error parsing folder ID: %v", err)
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	filePath, err := database.FilePath(folderId, userData.UserId)
	if err != nil {
		logger.LogError("Error retrieving file path: %v", err)
		http.Error(w, "Error retrieving file path", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(filePath))
}

func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("Items request received")
	userData, ok := r.Context().Value(middleware.UserDataKey).(database.UserData)
	logger.LogDebug("Items: User data retrieved from context: %v", userData)
	if !ok || userData.UserId == 0 {
		logger.LogError("User data not found or invalid in context")
		http.Error(w, "User data not found", http.StatusInternalServerError)
		return
	}

	folderId, err := extractFolderId(r)
	if err != nil {
		logger.LogError("Error parsing folder ID: %v", err)
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	files, err := database.GetFiles(folderId, userData.UserId)
	if err != nil {
		logger.LogError("Error retrieving files for user %d, folder %d: %v",
			userData.UserId, folderId, err)
		http.Error(w, "Unable to retrieve files", http.StatusInternalServerError)
		return
	}
	logger.LogDebug("Files retrieved successfully: %v", files)

	folders, err := database.GetFolders(folderId, userData.UserId)
	if err != nil {
		logger.LogError("Error retrieving folders for user %d, folder %d: %v",
			userData.UserId, folderId, err)
		http.Error(w, "Unable to retrieve folders", http.StatusInternalServerError)
		return
	}
	logger.LogDebug("Folders retrieved successfully: %v", folders)

	// Combine files and folders into a single slice of Items
	items := make([]models.Item, 0, len(files)+len(folders))
	for _, folder := range folders {
		items = append(items, folder)
	}
	for _, file := range files {
		items = append(items, file)
	}

	w.Header().Set("Content-Type", "text/html")
	err = components.TableComponent(items).Render(r.Context(), w)
	if err != nil {
		logger.LogError("Error rendering table: %v", err)
		http.Error(w, "Error rendering table", http.StatusInternalServerError)
		return
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("Upload request received")
	userData, ok := r.Context().Value(middleware.UserDataKey).(database.UserData)
	logger.LogDebug("Upload: user data retrieved from context: %v", userData)
	if !ok || userData.UserId == 0 {
		logger.LogError("User data not found or invalid in context")
		http.Error(w, "User data not found", http.StatusInternalServerError)
		return
	}

	folderIdStr := r.Header.Get("X-Folder-ID")
	folderId, err := strconv.ParseInt(folderIdStr, 10, 64)
	if err != nil {
		logger.LogError("Error parsing folder ID: %v", err)
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form (64MB max)
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		logger.LogError("Error parsing multipart form: %v", err)
		http.Error(w, "Error processing file upload", http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.LogError("Error retrieving file from form: %v", err)
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file contents
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.LogError("Error reading file: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Create file record in database
	fileData := models.UploadFile{
		UserId:    userData.UserId,
		FileName:  header.Filename,
		FolderId:  folderId,
		Content:   fileBytes,
		Size:      header.Size,
		CreatedAt: time.Now(),
	}

	// implementing upload function
	err = database.SaveFile(fileData)
	if err != nil {
		logger.LogError("Error saving file to database: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	logger.LogInfo("File uploaded successfully: %s", header.Filename)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("HX-Trigger", `{"upload" : "success"}`)
	w.WriteHeader(http.StatusOK)
}

func extractFolderId(r *http.Request) (int64, error) {
	folderIdStr := r.Header.Get("X-Folder-ID")
	if folderIdStr == "" {
		return 0, fmt.Errorf("missing folder ID")
	}

	folderId, err := strconv.ParseInt(folderIdStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid folder ID: %v", err)
	}

	return folderId, nil
}

// TODO work through this
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogInfo("Download request received")
	userData, ok := r.Context().Value(middleware.UserDataKey).(database.UserData)
	logger.LogDebug("Download: user data retrieved from context: %v", userData)
	if !ok || userData.UserId == 0 {
		logger.LogError("User data not found or invalid in context")
		http.Error(w, "User data not found", http.StatusInternalServerError)
		return
	}

	// Extract file ID from URL path (/download/{fileId})
	path := strings.TrimPrefix(r.URL.Path, "/download/")
	fileId, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		logger.LogError("Error parsing file ID from path: %v", err)
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	fileData, err := database.GetFile(fileId, userData.UserId)
	if err != nil {
		logger.LogError("Error retrieving file from database: %v", err)
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileData.FileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(fileData.Content)
	if err != nil {
		logger.LogError("Error writing file to response: %v", err)
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
}
