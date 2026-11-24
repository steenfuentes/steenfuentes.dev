package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/steenfuentes/portfolio-backend/internal/database"
	"github.com/steenfuentes/portfolio-backend/internal/models"
)

const (
	SessionCookieName = "portfolio_session"
	SessionDuration   = 24 * time.Hour * 7 // 7 days
)

type Manager struct {
	db *database.DB
}

func NewManager(db *database.DB) *Manager {
	return &Manager{db: db}
}

func (m *Manager) CreateSession(userID, githubUsername string) (*models.Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	session := &models.Session{
		ID:             sessionID,
		UserID:         userID,
		GitHubUsername: githubUsername,
		ExpiresAt:      time.Now().Add(SessionDuration),
		CreatedAt:      time.Now(),
	}

	query := `INSERT INTO sessions (id, user_id, github_username, expires_at) VALUES (?, ?, ?, ?)`
	_, err = m.db.Exec(query, session.ID, session.UserID, session.GitHubUsername, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (m *Manager) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	query := `SELECT id, user_id, github_username, expires_at, created_at FROM sessions WHERE id = ? AND expires_at > ?`
	
	err := m.db.QueryRow(query, sessionID, time.Now()).Scan(
		&session.ID, &session.UserID, &session.GitHubUsername, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (m *Manager) DeleteSession(sessionID string) error {
	_, err := m.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

func (m *Manager) CleanupExpiredSessions() error {
	_, err := m.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}

func (m *Manager) SetSessionCookie(w http.ResponseWriter, session *models.Session, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (m *Manager) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func (m *Manager) GetSessionFromRequest(r *http.Request) (*models.Session, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}

	return m.GetSession(cookie.Value)
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}