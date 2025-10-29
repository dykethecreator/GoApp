package domain

import (
	"time"

	"github.com/google/uuid"
)

// Poll represents a poll in a chat.
type Poll struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ChatID          uuid.UUID `json:"chat_id" db:"chat_id"`
	CreatedByUserID uuid.UUID `json:"created_by_user_id" db:"created_by_user_id"`
	QuestionText    string    `json:"question_text" db:"question_text"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// PollOption represents an option in a poll.
type PollOption struct {
	ID         int64     `json:"id" db:"id"`
	PollID     uuid.UUID `json:"poll_id" db:"poll_id"`
	OptionText string    `json:"option_text" db:"option_text"`
}

// PollVote represents a user's vote in a poll.
type PollVote struct {
	PollID   uuid.UUID `json:"poll_id" db:"poll_id"`
	OptionID int64     `json:"option_id" db:"option_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
}
