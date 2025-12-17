package models

import "time"

type UserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Username string   `json:"username"`
	Password string   `json:"password" binding:"required,min=3"`
	Roles    []string `json:"roles"`
}

type UserResponse struct {
	Email                 string    `json:"email" binding:"required,email"`
	Username              string    `json:"username"`
	RegistrationTimestamp time.Time `json:"registrationTimestamp"`
}
