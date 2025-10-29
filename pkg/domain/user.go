package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID                uuid.UUID `json:"id" db:"id"`
	PhoneNumber       string    `json:"phone_number" db:"phone_number"`
	DisplayName       string    `json:"display_name" db:"display_name"`
	ProfilePictureURL string    `json:"profile_picture_url" db:"profile_picture_url"`
	AboutText         string    `json:"about_text" db:"about_text"`
	LastSeenAt        time.Time `json:"last_seen_at" db:"last_seen_at"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// UserDevice represents a device a user has logged in with.
type UserDevice struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	UserID                uuid.UUID  `json:"user_id" db:"user_id"`
	RefreshTokenHash      string     `json:"-" db:"refresh_token_hash"`
	DeviceName            string     `json:"device_name" db:"device_name"`
	DeviceType            string     `json:"device_type" db:"device_type"`
	PushNotificationToken string     `json:"push_notification_token" db:"push_notification_token"`
	LastLoginAt           time.Time  `json:"last_login_at" db:"last_login_at"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	RevokedAt             *time.Time `json:"revoked_at" db:"revoked_at"`
}

// BlockedUser represents a blocked user relationship.
type BlockedUser struct {
	BlockerUserID uuid.UUID `json:"blocker_user_id" db:"blocker_user_id"`
	BlockedUserID uuid.UUID `json:"blocked_user_id" db:"blocked_user_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
