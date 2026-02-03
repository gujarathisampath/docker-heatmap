package models

import (
	"time"

	"gorm.io/gorm"
)

type EventType string

const (
	EventTypePush  EventType = "push"
	EventTypePull  EventType = "pull"
	EventTypeBuild EventType = "build"
)

type ActivityEvent struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign Key
	DockerAccountID uint          `gorm:"column:docker_account_id;not null;index" json:"docker_account_id"`
	DockerAccount   DockerAccount `gorm:"foreignKey:DockerAccountID" json:"-"`

	// Event Data
	EventType EventType `gorm:"column:event_type;not null;index" json:"event_type"`
	EventDate time.Time `gorm:"column:event_date;not null;index" json:"event_date"`
	Count     int       `gorm:"column:count;not null;default:1" json:"count"`

	// Repository Info
	Repository string `gorm:"column:repository" json:"repository,omitempty"`
	Tag        string `gorm:"column:tag" json:"tag,omitempty"`
}

// TableName specifies the table name
func (ActivityEvent) TableName() string {
	return "activity_events"
}

func (a *ActivityEvent) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	// Normalize event date to midnight UTC
	a.EventDate = time.Date(
		a.EventDate.Year(),
		a.EventDate.Month(),
		a.EventDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
	return nil
}

// ActivitySummary represents aggregated activity for a specific date
type ActivitySummary struct {
	Date       string `json:"date"`
	TotalCount int    `json:"count"`
	Pushes     int    `json:"pushes"`
	Pulls      int    `json:"pulls"`
	Builds     int    `json:"builds"`
	Level      int    `json:"level"`
}
