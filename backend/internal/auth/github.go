package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/steenfuentes/portfolio-backend/internal/database"
	"github.com/steenfuentes/portfolio-backend/internal/models"
	"github.com/steenfuentes/portfolio-backend/internal/session"
)

type GitHubAuth struct {
	clientID       string
	clientSecret   string
	redirectURI    string
	allowedUser    string
	sessionManager *session.Manager
}

func NewGitHubAuth(db *database.DB) *GitHubAuth {
	return &GitHubAuth{
		clientID:       os.Getenv("GITHUB_CLIENT_ID"),
		clientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		redirectURI:    os.Getenv("GITHUB_REDIRECT_URI"),
		allowedUser:    os.Getenv("GITHUB_ALLOWED_USER"),
		sessionManager: session.NewManager(db),
	}
}

func (g *GitHubAuth) GetAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", g.clientID)
	params.Add("redirect_uri", g.redirectURI)
	params.Add("scope", "read:user user:email")
	params.Add("state", state)

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params.Encode())
}

func (g *GitHubAuth) ExchangeCode(ctx context.Context, code string) (string, error) {
	params := url.Values{}
	params.Add("client_id", g.clientID)
	params.Add("client_secret", g.clientSecret)
	params.Add("code", code)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return "", err
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("github oauth error: %s", result.Error)
	}

	return result.AccessToken, nil
}

func (g *GitHubAuth) GetUser(ctx context.Context, accessToken string) (*models.GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (g *GitHubAuth) IsUserAllowed(username string) bool {
	return g.allowedUser == "" || g.allowedUser == username
}

func (g *GitHubAuth) CreateSession(user *models.GitHubUser) (*models.Session, error) {
	return g.sessionManager.CreateSession(strconv.FormatInt(user.ID, 10), user.Login)
}

func (g *GitHubAuth) GenerateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (g *GitHubAuth) SessionManager() *session.Manager {
	return g.sessionManager
}