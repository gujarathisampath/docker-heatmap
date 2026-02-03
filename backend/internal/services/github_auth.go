package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"docker-heatmap/internal/config"
	"docker-heatmap/internal/database"
	"docker-heatmap/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	ErrGitHubAuthFailed = errors.New("github authentication failed")
	ErrUserNotFound     = errors.New("user not found")
)

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type GitHubAuthService struct {
	oauthConfig *oauth2.Config
}

func NewGitHubAuthService() *GitHubAuthService {
	return &GitHubAuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     config.AppConfig.GitHubClientID,
			ClientSecret: config.AppConfig.GitHubClientSecret,
			RedirectURL:  config.AppConfig.GitHubCallbackURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// GetAuthURL returns the GitHub OAuth authorization URL
func (s *GitHubAuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeCode exchanges the authorization code for access token and fetches user data
func (s *GitHubAuthService) ExchangeCode(ctx context.Context, code string) (*models.User, error) {
	// Exchange code for token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGitHubAuthFailed, err)
	}

	// Fetch user data from GitHub
	githubUser, err := s.fetchGitHubUser(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	// Find or create user in database
	user, err := s.findOrCreateUser(githubUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *GitHubAuthService) fetchGitHubUser(ctx context.Context, accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrGitHubAuthFailed
	}

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	// Fetch primary email if not public
	if githubUser.Email == "" {
		email, _ := s.fetchPrimaryEmail(ctx, accessToken)
		githubUser.Email = email
	}

	return &githubUser, nil
}

func (s *GitHubAuthService) fetchPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary {
			return e.Email, nil
		}
	}

	return "", nil
}

func (s *GitHubAuthService) findOrCreateUser(githubUser *GitHubUser) (*models.User, error) {
	var user models.User

	// Try to find existing user
	result := database.DB.Where("github_id = ?", githubUser.ID).First(&user)
	if result.Error == nil {
		// Update user data
		user.GitHubUsername = githubUser.Login
		user.GitHubEmail = githubUser.Email
		user.AvatarURL = githubUser.AvatarURL
		user.Name = githubUser.Name
		database.DB.Save(&user)
		return &user, nil
	}

	// Create new user
	user = models.User{
		GitHubID:       githubUser.ID,
		GitHubUsername: githubUser.Login,
		GitHubEmail:    githubUser.Email,
		AvatarURL:      githubUser.AvatarURL,
		Name:           githubUser.Name,
		PublicProfile:  true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID fetches a user by their ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// GetUserByGitHubUsername fetches a user by their GitHub username
func GetUserByGitHubUsername(username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("github_username = ?", username).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}
