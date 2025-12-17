package models

import "time"

type UserBoundary struct {
	Email                 string    `json:"email" binding:"required,email"`
	Username              string    `json:"username"`
	Password              string    `json:"password" binding:"required,min=3"`
	Roles                 []string  `json:"roles"`
	RegistrationTimestamp time.Time `json:"registrationTimestamp"`
}

type UserEntity struct {
	ID                    string
	Email                 string
	Username              string
	Password              string
	Roles                 []string
	RegistrationTimestamp time.Time
}
