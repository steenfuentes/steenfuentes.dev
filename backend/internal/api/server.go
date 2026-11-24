package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/steenfuentes/portfolio-backend/internal/auth"
	"github.com/steenfuentes/portfolio-backend/internal/database"
	"github.com/steenfuentes/portfolio-backend/internal/models"
)

type Server struct {
	db          *database.DB
	auth        *auth.GitHubAuth
	mux         *http.ServeMux
	states      map[string]bool // Simple in-memory state store
	authEnabled bool
}

func NewServer(db *database.DB) *Server {
	authEnabled := os.Getenv("AUTH_ENABLED") != "false"
	
	s := &Server{
		db:          db,
		auth:        auth.NewGitHubAuth(db),
		mux:         http.NewServeMux(),
		states:      make(map[string]bool),
		authEnabled: authEnabled,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Health check
	s.mux.HandleFunc("GET /health", s.handleHealth)
	
	// Public routes
	s.mux.HandleFunc("GET /api/posts", s.handleGetPosts)
	s.mux.HandleFunc("GET /api/posts/{slug}", s.handleGetPost)

	// Auth routes
	s.mux.HandleFunc("GET /api/auth/github", s.handleGitHubAuth)
	s.mux.HandleFunc("GET /api/auth/github/callback", s.handleGitHubCallback)
	s.mux.HandleFunc("GET /api/auth/me", s.handleGetMe)
	s.mux.HandleFunc("POST /api/auth/logout", s.handleLogout)

	// Protected routes
	s.mux.HandleFunc("POST /api/posts", s.requireAuth(s.handleCreatePost))
	s.mux.HandleFunc("PUT /api/posts/{id}", s.requireAuth(s.handleUpdatePost))
	s.mux.HandleFunc("DELETE /api/posts/{id}", s.requireAuth(s.handleDeletePost))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(s.mux).ServeHTTP(w, r)
}

// Health check handler
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "healthy",
		"service": "backend",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"database": "connected",
	}
	
	// Check database connection
	if err := s.db.Ping(); err != nil {
		health["status"] = "unhealthy"
		health["database"] = "disconnected"
		s.respondJSON(w, http.StatusServiceUnavailable, health)
		return
	}
	
	s.respondJSON(w, http.StatusOK, health)
}

// getAllowedOrigins returns the list of allowed CORS origins based on environment
func (s *Server) getAllowedOrigins() []string {
	env := os.Getenv("ENVIRONMENT")
	frontendURL := os.Getenv("FRONTEND_URL")
	
	if env == "development" {
		return []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:4173"}
	}
	
	if env == "preview" {
		// Allow all origins for preview environments
		return []string{"*"}
	}
	
	// Production: only allow specific frontend URL
	origins := []string{}
	if frontendURL != "" {
		origins = append(origins, frontendURL)
		// Also allow the www version if present
		if frontendURL == "https://steenfuentes.dev" {
			origins = append(origins, "https://www.steenfuentes.dev")
		}
	}
	
	return origins
}

// Posts handlers
func (s *Server) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10
	publishedOnly := true

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Check if user is authenticated to see unpublished posts
	if session, _ := s.auth.SessionManager().GetSessionFromRequest(r); session != nil {
		publishedOnly = false
	}

	posts, err := s.db.GetPosts(page, pageSize, publishedOnly)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	s.respondJSON(w, http.StatusOK, posts)
}

func (s *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := s.db.GetPostBySlug(slug)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch post")
		return
	}

	if post == nil {
		s.respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Check if post is published or user is authenticated
	if !post.Published {
		if session, _ := s.auth.SessionManager().GetSessionFromRequest(r); session == nil {
			s.respondError(w, http.StatusNotFound, "Post not found")
			return
		}
	}

	s.respondJSON(w, http.StatusOK, post)
}

func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	post, err := s.db.CreatePost(&req)
	if err != nil {
		log.Printf("Failed to create post: %v", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	s.respondJSON(w, http.StatusCreated, post)
}

func (s *Server) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req models.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	post, err := s.db.UpdatePost(id, &req)
	if err != nil {
		log.Printf("Failed to update post: %v", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	if post == nil {
		s.respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	s.respondJSON(w, http.StatusOK, post)
}

func (s *Server) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.db.DeletePost(id); err != nil {
		log.Printf("Failed to delete post: %v", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Auth handlers
func (s *Server) handleGitHubAuth(w http.ResponseWriter, r *http.Request) {
	// If auth is disabled, return an error
	if !s.authEnabled {
		s.respondError(w, http.StatusForbidden, "Authentication is disabled")
		return
	}

	state, err := s.auth.GenerateState()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to generate state")
		return
	}

	s.states[state] = true
	authURL := s.auth.GetAuthURL(state)
	
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (s *Server) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// If auth is disabled, return an error
	if !s.authEnabled {
		s.respondError(w, http.StatusForbidden, "Authentication is disabled")
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		s.respondError(w, http.StatusBadRequest, "Missing code or state")
		return
	}

	// Validate state
	if !s.states[state] {
		s.respondError(w, http.StatusBadRequest, "Invalid state")
		return
	}
	delete(s.states, state)

	// Exchange code for access token
	accessToken, err := s.auth.ExchangeCode(r.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	// Get user info
	user, err := s.auth.GetUser(r.Context(), accessToken)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to get user info")
		return
	}

	// Check if user is allowed
	if !s.auth.IsUserAllowed(user.Login) {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Create session
	session, err := s.auth.CreateSession(user)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	// Set session cookie
	secure := os.Getenv("ENVIRONMENT") == "production"
	s.auth.SessionManager().SetSessionCookie(w, session, secure)

	// Redirect to frontend
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	http.Redirect(w, r, fmt.Sprintf("%s/admin", frontendURL), http.StatusTemporaryRedirect)
}

func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) {
	// If auth is disabled, return a mock user
	if !s.authEnabled {
		s.respondJSON(w, http.StatusOK, models.AuthResponse{
			User: models.GitHubUser{
				Login: "dev-user",
			},
		})
		return
	}

	session, err := s.auth.SessionManager().GetSessionFromRequest(r)
	if err != nil || session == nil {
		s.respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	s.respondJSON(w, http.StatusOK, models.AuthResponse{
		User: models.GitHubUser{
			Login: session.GitHubUsername,
		},
		ExpiresAt: session.ExpiresAt,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.auth.SessionManager().GetSessionFromRequest(r)
	if session != nil {
		s.auth.SessionManager().DeleteSession(session.ID)
	}

	s.auth.SessionManager().ClearSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}