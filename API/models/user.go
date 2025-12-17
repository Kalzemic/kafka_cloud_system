package models

import "time"

type User struct {
	Email                 string    `json:"email" binding:"required,email"`
	Username              string    `json:"username"`
	Password              string    `json:"password" binding:"required,min=3"`
	Roles                 []string  `json:"roles"`
	RegistrationTimestamp time.Time `json:"registrationTimestamp"`
}
