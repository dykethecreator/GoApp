package domain

import (
	"time"

	"github.com/google/uuid"
)

// StatusContentType defines the content type of a status update.
type StatusContentType string

const (
	TextStatus  StatusContentType = "text"
	ImageStatus StatusContentType = "image"
	VideoStatus StatusContentType = "video"
)

// StatusUpdate represents a user's status update (story).
type StatusUpdate struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	UserID       uuid.UUID         `json:"user_id" db:"user_id"`
	ContentType  StatusContentType `json:"content_type" db:"content_type"`
	ContentOrURL string            `json:"content_or_url" db:"content_or_url"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	ExpiresAt    time.Time         `json:"expires_at" db:"expires_at"`
}

// StatusView represents a view of a status update by a user.
type StatusView struct {
	StatusID uuid.UUID `json:"status_id" db:"status_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	ViewedAt time.Time `json:"viewed_at" db:"viewed_at"`
}
