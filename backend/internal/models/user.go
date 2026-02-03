package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// GitHub OAuth Data
	GitHubID       int64  `gorm:"column:github_id;uniqueIndex;not null" json:"github_id"`
	GitHubUsername string `gorm:"column:github_username;not null" json:"github_username"`
	GitHubEmail    string `gorm:"column:github_email" json:"email,omitempty"`
	AvatarURL      string `gorm:"column:avatar_url" json:"avatar_url,omitempty"`
	Name           string `gorm:"column:name" json:"name,omitempty"`

	// Profile Settings
	PublicProfile bool   `gorm:"column:public_profile;default:true" json:"public_profile"`
	Bio           string `gorm:"column:bio" json:"bio,omitempty"`

	// Relationships
	DockerAccounts []DockerAccount `gorm:"foreignKey:UserID" json:"docker_accounts,omitempty"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
