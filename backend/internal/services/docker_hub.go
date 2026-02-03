package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"docker-heatmap/internal/config"
	"docker-heatmap/internal/database"
	"docker-heatmap/internal/models"
	"docker-heatmap/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrDockerHubAuthFailed   = errors.New("docker hub authentication failed")
	ErrDockerAccountNotFound = errors.New("docker account not found")
	ErrDockerAccountExists   = errors.New("docker account already connected")
	ErrInvalidDockerToken    = errors.New("invalid docker hub access token")
)

// Use shared HTTP client from utils package
var httpClient = utils.HTTPClient

// parseDockerTime parses Docker Hub's date format which includes microseconds
func parseDockerTime(dateStr string) (time.Time, error) {
	// Docker Hub uses ISO 8601 format with microseconds: 2026-01-17T08:19:30.340959Z
	formats := []string{
		"2006-01-02T15:04:05.999999Z",
		"2006-01-02T15:04:05.999999-07:00",
		time.RFC3339Nano,
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

type DockerHubService struct {
	apiURL string
}

func NewDockerHubService() *DockerHubService {
	return &DockerHubService{
		apiURL: config.AppConfig.DockerHubAPIURL,
	}
}

// DockerHubRepository represents a repository from Docker Hub API
type DockerHubRepository struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Description string `json:"description"`
	LastUpdated string `json:"last_updated"` // ISO string
	PullCount   int64  `json:"pull_count"`
	StarCount   int    `json:"star_count"`
	IsPrivate   bool   `json:"is_private"`
}

// DockerHubTag represents a tag from Docker Hub API
type DockerHubTag struct {
	Name          string `json:"name"`
	LastUpdated   string `json:"last_updated"`
	TagLastPushed string `json:"tag_last_pushed"`
	Digest        string `json:"digest"`
}

// ConnectAccount validates and connects a Docker Hub account.
// It aggressively handles duplicates by cleaning up any previous records for the user or username.
func (s *DockerHubService) ConnectAccount(ctx context.Context, userID uint, dockerUsername, accessToken string) (*models.DockerAccount, error) {
	var account models.DockerAccount

	// Use a transaction to ensure deletion and creation are atomic
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Check for username conflict with OTHER users (including soft-deleted)
		var conflict models.DockerAccount
		if err := tx.Unscoped().Where("docker_username = ? AND user_id != ?", dockerUsername, userID).First(&conflict).Error; err == nil {
			return errors.New("this Docker username is connected to another account")
		}

		// 2. Clear any existing/orphaned records and their events
		var accountIDs []uint
		tx.Unscoped().Model(&models.DockerAccount{}).Where("user_id = ? OR docker_username = ?", userID, dockerUsername).Pluck("id", &accountIDs)

		if len(accountIDs) > 0 {
			log.Printf("ConnectAccount: Found %d existing records to delete (IDs: %v)", len(accountIDs), accountIDs)
			// Delete events first to satisfy foreign key constraints
			if err := tx.Unscoped().Where("docker_account_id IN ?", accountIDs).Delete(&models.ActivityEvent{}).Error; err != nil {
				return fmt.Errorf("failed to clear old events: %w", err)
			}
			// Delete the accounts
			if err := tx.Unscoped().Where("id IN ?", accountIDs).Delete(&models.DockerAccount{}).Error; err != nil {
				return fmt.Errorf("failed to clear old account: %w", err)
			}
			log.Printf("ConnectAccount: Successfully cleared records for UserID=%d and DockerUser=%s", userID, dockerUsername)
		}

		// 3. Validate the Docker Hub username exists
		if err := s.validateUsername(ctx, dockerUsername); err != nil {
			return err
		}

		// 4. Encrypt the access token
		encryptedToken, iv, err := utils.Encrypt(accessToken)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		// 5. Create a fresh account record
		account = models.DockerAccount{
			UserID:         userID,
			DockerUsername: dockerUsername,
			EncryptedToken: encryptedToken,
			TokenIV:        iv,
			IsActive:       true,
			AutoRefresh:    true,
		}

		if err := tx.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("ConnectAccount failed: %v", err)
		return nil, err
	}

	log.Printf("Docker account connected: ID=%d, Username=%s for User=%d", account.ID, dockerUsername, userID)

	// 6. Trigger initial sync
	go func() {
		syncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := s.SyncActivity(syncCtx, account.ID); err != nil {
			log.Printf("Initial sync failed for account %d: %v", account.ID, err)
		}
	}()

	return &account, nil
}

// validateUsername checks if a Docker Hub username exists
func (s *DockerHubService) validateUsername(ctx context.Context, username string) error {
	url := fmt.Sprintf("%s/users/%s", s.apiURL, username)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Docker Hub API request failed: %v", err)
		return ErrDockerHubAuthFailed
	}
	defer resp.Body.Close()

	log.Printf("Docker Hub user lookup: status=%d for user=%s", resp.StatusCode, username)

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("Docker Hub user not found")
	}

	if resp.StatusCode != http.StatusOK {
		return ErrDockerHubAuthFailed
	}

	return nil
}

// GetDockerAccount retrieves a user's Docker account
func (s *DockerHubService) GetDockerAccount(userID uint) (*models.DockerAccount, error) {
	var account models.DockerAccount
	if err := database.DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		return nil, ErrDockerAccountNotFound
	}
	return &account, nil
}

// GetDockerAccountByUsername retrieves a Docker account by Docker username
func (s *DockerHubService) GetDockerAccountByUsername(dockerUsername string) (*models.DockerAccount, error) {
	var account models.DockerAccount
	if err := database.DB.Where("docker_username = ?", dockerUsername).First(&account).Error; err != nil {
		return nil, ErrDockerAccountNotFound
	}
	return &account, nil
}

// FetchRepositories fetches repositories for a Docker Hub user
func (s *DockerHubService) FetchRepositories(ctx context.Context, username, token string) ([]DockerHubRepository, error) {
	url := fmt.Sprintf("%s/repositories/%s?page_size=100", s.apiURL, username)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	client := httpClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Docker Hub API error: status=%d body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("docker hub api returned status %d", resp.StatusCode)
	}

	var result struct {
		Count   int                   `json:"count"`
		Results []DockerHubRepository `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Fetched %d repositories for %s", len(result.Results), username)
	return result.Results, nil
}

// FetchTags fetches tags for a specific repository
func (s *DockerHubService) FetchTags(ctx context.Context, username, repoName, token string) ([]DockerHubTag, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/tags?page_size=100", s.apiURL, username, repoName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker hub api returned status %d", resp.StatusCode)
	}

	var result struct {
		Results []DockerHubTag `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// SyncActivity syncs Docker Hub activity for an account
func (s *DockerHubService) SyncActivity(ctx context.Context, accountID uint) error {
	var account models.DockerAccount
	if err := database.DB.First(&account, accountID).Error; err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	log.Printf("Starting sync for account ID=%d username=%s", account.ID, account.DockerUsername)

	// Mark sync in progress
	account.SyncInProgress = true
	database.DB.Save(&account)

	defer func() {
		account.SyncInProgress = false
		now := time.Now()
		account.LastSyncAt = &now
		database.DB.Save(&account)
		log.Printf("Sync completed for account ID=%d", account.ID)
	}()

	// Decrypt token
	token, err := utils.Decrypt(account.EncryptedToken, account.TokenIV)
	if err != nil {
		log.Printf("Failed to decrypt token for account %d: %v", account.ID, err)
		// Continue anyway, maybe public access works
		token = ""
	}

	// Fetch repositories
	repos, err := s.FetchRepositories(ctx, account.DockerUsername, token)
	if err != nil {
		account.LastSyncError = err.Error()
		log.Printf("Failed to fetch repositories for %s: %v", account.DockerUsername, err)
		return err
	}

	if len(repos) == 0 {
		log.Printf("No repositories found for %s", account.DockerUsername)
		account.LastSyncError = ""
		return nil
	}

	log.Printf("Processing %d repositories for %s", len(repos), account.DockerUsername)

	eventsCreated := 0

	// Process each repository to extract activity
	for _, repo := range repos {
		log.Printf("Processing repo: %s, last_updated: %s", repo.Name, repo.LastUpdated)

		// Parse and create push event based on last_updated
		if repo.LastUpdated != "" {
			parsedTime, err := parseDockerTime(repo.LastUpdated)
			if err != nil {
				log.Printf("Failed to parse date %s: %v", repo.LastUpdated, err)
			} else if s.createActivity(&account, models.EventTypePush, parsedTime, repo.Name, "") {
				eventsCreated++
			}
		}

		// Fetch tags and create push events for each
		tags, err := s.FetchTags(ctx, account.DockerUsername, repo.Name, token)
		if err != nil {
			log.Printf("Failed to fetch tags for %s/%s: %v", account.DockerUsername, repo.Name, err)
			continue
		}

		for _, tag := range tags {
			// Use tag_last_pushed
			if tag.TagLastPushed != "" {
				parsedTime, err := parseDockerTime(tag.TagLastPushed)
				if err != nil {
					log.Printf("Failed to parse tag date %s: %v", tag.TagLastPushed, err)
				} else if s.createActivity(&account, models.EventTypePush, parsedTime, repo.Name, tag.Name) {
					eventsCreated++
				}
			}
		}
	}

	log.Printf("Created/updated %d activity events for %s", eventsCreated, account.DockerUsername)
	account.LastSyncError = ""
	return nil
}

// createActivity creates an activity event (returns true if created)
func (s *DockerHubService) createActivity(account *models.DockerAccount, eventType models.EventType, eventDate time.Time, repo, tag string) bool {
	// Normalize to date only (midnight UTC)
	normalizedDate := time.Date(
		eventDate.Year(),
		eventDate.Month(),
		eventDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	// Check if event already exists
	var existing models.ActivityEvent
	err := database.DB.Where(
		"docker_account_id = ? AND event_date = ? AND repository = ? AND tag = ?",
		account.ID, normalizedDate, repo, tag,
	).First(&existing).Error

	if err == nil {
		// Event already exists, update count
		existing.Count++
		database.DB.Save(&existing)
		return false
	}

	// Create new event
	event := models.ActivityEvent{
		DockerAccountID: account.ID,
		EventType:       eventType,
		EventDate:       normalizedDate,
		Repository:      repo,
		Tag:             tag,
		Count:           1,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		log.Printf("Failed to create activity event: %v", err)
		return false
	}

	return true
}

// GetActivitySummary returns aggregated activity data for heatmap
func (s *DockerHubService) GetActivitySummary(dockerUsername string, days int) ([]models.ActivitySummary, error) {
	account, err := s.GetDockerAccountByUsername(dockerUsername)
	if err != nil {
		return nil, err
	}

	startDate := time.Now().UTC().AddDate(0, 0, -days)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Now().UTC()

	// Query activity events
	var events []models.ActivityEvent
	err = database.DB.Where(
		"docker_account_id = ? AND event_date >= ? AND event_date <= ?",
		account.ID, startDate, endDate,
	).Find(&events).Error

	if err != nil {
		log.Printf("Failed to query activity events: %v", err)
		return nil, err
	}

	log.Printf("Found %d activity events for %s from %s to %s", len(events), dockerUsername, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Aggregate by date
	dateMap := make(map[string]*models.ActivitySummary)

	for _, event := range events {
		dateStr := event.EventDate.Format("2006-01-02")
		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = &models.ActivitySummary{Date: dateStr}
		}

		summary := dateMap[dateStr]
		summary.TotalCount += event.Count

		switch event.EventType {
		case models.EventTypePush:
			summary.Pushes += event.Count
		case models.EventTypePull:
			summary.Pulls += event.Count
		case models.EventTypeBuild:
			summary.Builds += event.Count
		}
	}

	// Calculate max for level calculation
	maxCount := 0
	for _, s := range dateMap {
		if s.TotalCount > maxCount {
			maxCount = s.TotalCount
		}
	}

	// Fill all dates in range
	summaries := make([]models.ActivitySummary, 0, days+1)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if summary, exists := dateMap[dateStr]; exists {
			summary.Level = calculateLevel(summary.TotalCount, maxCount)
			summaries = append(summaries, *summary)
		} else {
			summaries = append(summaries, models.ActivitySummary{
				Date:  dateStr,
				Level: 0,
			})
		}
	}

	return summaries, nil
}

// calculateLevel calculates the intensity level (0-4) for heatmap
func calculateLevel(count, maxCount int) int {
	if count == 0 || maxCount == 0 {
		return 0
	}

	ratio := float64(count) / float64(maxCount)
	switch {
	case ratio > 0.75:
		return 4
	case ratio > 0.5:
		return 3
	case ratio > 0.25:
		return 2
	case ratio > 0:
		return 1
	default:
		return 0
	}
}

// DisconnectAccount removes a Docker Hub account permanently
func (s *DockerHubService) DisconnectAccount(userID, accountID uint) error {
	// Permanently delete all activity events (use Unscoped to bypass soft delete)
	database.DB.Unscoped().Where("docker_account_id = ?", accountID).Delete(&models.ActivityEvent{})

	// Permanently delete the docker account (use Unscoped to bypass soft delete)
	result := database.DB.Unscoped().Where("id = ? AND user_id = ?", accountID, userID).Delete(&models.DockerAccount{})
	if result.RowsAffected == 0 {
		return ErrDockerAccountNotFound
	}

	log.Printf("Permanently deleted docker account ID=%d for user ID=%d", accountID, userID)
	return nil
}
