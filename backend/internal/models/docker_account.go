package models

import (
	"time"

	"gorm.io/gorm"
)

type DockerAccount struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign Key
	UserID uint `gorm:"column:user_id;not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	// Docker Hub Data
	DockerUsername string `gorm:"column:docker_username;not null;uniqueIndex" json:"docker_username"`

	// Encrypted Access Token (AES-256 encrypted)
	EncryptedToken string `gorm:"column:encrypted_token;not null" json:"-"`
	TokenIV        string `gorm:"column:token_iv;not null" json:"-"`

	// Sync Status
	LastSyncAt     *time.Time `gorm:"column:last_sync_at" json:"last_sync_at,omitempty"`
	LastSyncError  string     `gorm:"column:last_sync_error" json:"last_sync_error,omitempty"`
	SyncInProgress bool       `gorm:"column:sync_in_progress;default:false" json:"sync_in_progress"`

	// Settings
	IsActive    bool `gorm:"column:is_active;default:true" json:"is_active"`
	AutoRefresh bool `gorm:"column:auto_refresh;default:true" json:"auto_refresh"`

	// Relationships
	ActivityEvents []ActivityEvent `gorm:"foreignKey:DockerAccountID" json:"activity_events,omitempty"`
}

// TableName specifies the table name
func (DockerAccount) TableName() string {
	return "docker_accounts"
}

func (d *DockerAccount) BeforeCreate(tx *gorm.DB) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

func (d *DockerAccount) BeforeUpdate(tx *gorm.DB) error {
	d.UpdatedAt = time.Now()
	return nil
}
