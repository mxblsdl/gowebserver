package models

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}