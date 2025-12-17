package models

import "time"

type Post struct {
	UserEmail string    `json:"email"`
	Content   string    `json:"content"`
	TimeStamp time.Time `json:"timestamp"`
}
