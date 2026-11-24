package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/steenfuentes/portfolio-backend/internal/models"
)

func (db *DB) GetPosts(page, pageSize int, publishedOnly bool) (*models.PostsResponse, error) {
	offset := (page - 1) * pageSize
	
	query := `SELECT id, slug, title, content, excerpt, tags, published, created_at, updated_at, reading_time 
			  FROM posts`
	countQuery := `SELECT COUNT(*) FROM posts`
	
	if publishedOnly {
		query += " WHERE published = true"
		countQuery += " WHERE published = true"
	}
	
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	
	var totalCount int
	err := db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}
	
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()
	
	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Content, &p.Excerpt, &p.Tags, &p.Published, &p.CreatedAt, &p.UpdatedAt, &p.ReadingTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, p)
	}
	
	return &models.PostsResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (db *DB) GetPostBySlug(slug string) (*models.Post, error) {
	var p models.Post
	query := `SELECT id, slug, title, content, excerpt, tags, published, created_at, updated_at, reading_time 
			  FROM posts WHERE slug = ?`
	
	err := db.QueryRow(query, slug).Scan(&p.ID, &p.Slug, &p.Title, &p.Content, &p.Excerpt, &p.Tags, &p.Published, &p.CreatedAt, &p.UpdatedAt, &p.ReadingTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	
	return &p, nil
}

func (db *DB) CreatePost(req *models.CreatePostRequest) (*models.Post, error) {
	id := uuid.New().String()
	postSlug := slug.Make(req.Title)
	now := time.Now()
	
	tagsJSON, _ := json.Marshal(req.Tags)
	readingTime := calculateReadingTime(req.Content)
	
	query := `INSERT INTO posts (id, slug, title, content, excerpt, tags, published, created_at, updated_at, reading_time)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, id, postSlug, req.Title, req.Content, req.Excerpt, string(tagsJSON), req.Published, now, now, readingTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	
	return &models.Post{
		ID:          id,
		Slug:        postSlug,
		Title:       req.Title,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		Tags:        string(tagsJSON),
		Published:   req.Published,
		CreatedAt:   now,
		UpdatedAt:   now,
		ReadingTime: readingTime,
	}, nil
}

func (db *DB) UpdatePost(id string, req *models.UpdatePostRequest) (*models.Post, error) {
	var updates []string
	var args []interface{}
	
	if req.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *req.Title)
		
		updates = append(updates, "slug = ?")
		args = append(args, slug.Make(*req.Title))
	}
	
	if req.Content != nil {
		updates = append(updates, "content = ?")
		args = append(args, *req.Content)
		
		updates = append(updates, "reading_time = ?")
		args = append(args, calculateReadingTime(*req.Content))
	}
	
	if req.Excerpt != nil {
		updates = append(updates, "excerpt = ?")
		args = append(args, *req.Excerpt)
	}
	
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		updates = append(updates, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	
	if req.Published != nil {
		updates = append(updates, "published = ?")
		args = append(args, *req.Published)
	}
	
	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())
	
	args = append(args, id)
	
	query := fmt.Sprintf("UPDATE posts SET %s WHERE id = ?", strings.Join(updates, ", "))
	
	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	
	return db.GetPostByID(id)
}

func (db *DB) DeletePost(id string) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (db *DB) GetPostByID(id string) (*models.Post, error) {
	var p models.Post
	query := `SELECT id, slug, title, content, excerpt, tags, published, created_at, updated_at, reading_time 
			  FROM posts WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Slug, &p.Title, &p.Content, &p.Excerpt, &p.Tags, &p.Published, &p.CreatedAt, &p.UpdatedAt, &p.ReadingTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	
	return &p, nil
}

func calculateReadingTime(content string) int {
	words := strings.Fields(content)
	wordsPerMinute := 200
	readingTime := len(words) / wordsPerMinute
	if readingTime < 1 {
		readingTime = 1
	}
	return readingTime
}