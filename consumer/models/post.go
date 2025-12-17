package models

import "time"

type PollRequest struct {
	MaxPosts    int           `json:"maxPosts"`
	MaxDuration time.Duration `json:"maxDuration"`
}

type Post struct {
	UserEmail string    `json:"email"`
	Content   string    `json:"content"`
	TimeStamp time.Time `json:"timestamp"`
}
