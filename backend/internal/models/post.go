package models

import (
	"time"
)

type Post struct {
	ID          string    `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Excerpt     string    `json:"excerpt" db:"excerpt"`
	Tags        string    `json:"tags" db:"tags"` // JSON array stored as string
	Published   bool      `json:"published" db:"published"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	ReadingTime int       `json:"reading_time" db:"reading_time"`
}

type CreatePostRequest struct {
	Title     string   `json:"title" validate:"required,min=1,max=200"`
	Content   string   `json:"content" validate:"required"`
	Excerpt   string   `json:"excerpt" validate:"max=500"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

type UpdatePostRequest struct {
	Title     *string  `json:"title" validate:"omitempty,min=1,max=200"`
	Content   *string  `json:"content" validate:"omitempty"`
	Excerpt   *string  `json:"excerpt" validate:"omitempty,max=500"`
	Tags      []string `json:"tags"`
	Published *bool    `json:"published"`
}

type PostsResponse struct {
	Posts      []Post `json:"posts"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}