package domain

import (
	"time"

	"github.com/google/uuid"
)

// CallType defines the type of a call.
type CallType string

const (
	AudioCall CallType = "audio"
	VideoCall CallType = "video"
)

// CallStatus defines the status of a call.
type CallStatus string

const (
	CompletedCall CallStatus = "completed"
	MissedCall    CallStatus = "missed"
	RejectedCall  CallStatus = "rejected"
)

// CallLog represents a record of a call.
type CallLog struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ChatID         uuid.UUID  `json:"chat_id" db:"chat_id"`
	CallerUserID   uuid.UUID  `json:"caller_user_id" db:"caller_user_id"`
	CallType       CallType   `json:"call_type" db:"call_type"`
	Status         CallStatus `json:"status" db:"status"`
	StartedAt      time.Time  `json:"started_at" db:"started_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty" db:"ended_at"`
}
