package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/steenfuentes/portfolio-backend/internal/models"
)

type contextKey string

const (
	sessionContextKey contextKey = "session"
)

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If auth is disabled, create a mock session for development
		if !s.authEnabled {
			mockSession := &models.Session{
				ID:             "dev-session",
				UserID:         "dev-user",
				GitHubUsername: "dev-user",
				ExpiresAt:      time.Now().Add(24 * time.Hour),
			}
			ctx := context.WithValue(r.Context(), sessionContextKey, mockSession)
			next(w, r.WithContext(ctx))
			return
		}

		session, err := s.auth.SessionManager().GetSessionFromRequest(r)
		if err != nil || session == nil {
			s.respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), sessionContextKey, session)
		next(w, r.WithContext(ctx))
	}
}

func getSession(ctx context.Context) *models.Session {
	session, _ := ctx.Value(sessionContextKey).(*models.Session)
	return session
}

func (s *Server) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Allow specific origins based on environment
		allowedOrigins := s.getAllowedOrigins()
		isAllowed := false
		
		for _, allowed := range allowedOrigins {
			if origin == allowed || allowed == "*" {
				isAllowed = true
				break
			}
		}
		
		if isAllowed && origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if len(allowedOrigins) > 0 && allowedOrigins[0] == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, map[string]string{"error": message})
}