package models

import (
	"time"
)

type Session struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	GitHubUsername string    `json:"github_username" db:"github_username"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type AuthResponse struct {
	User      GitHubUser `json:"user"`
	ExpiresAt time.Time  `json:"expires_at"`
}