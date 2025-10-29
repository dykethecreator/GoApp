package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ContentType defines the type of message content.
type ContentType string

const (
	TextContent               ContentType = "text"
	ImageContent              ContentType = "image"
	VideoContent              ContentType = "video"
	AudioContent              ContentType = "audio"
	FileContent               ContentType = "file"
	SystemNotificationContent ContentType = "system_notification"
	DeletedContent            ContentType = "deleted"
	PollContent               ContentType = "poll"
)

// Message represents a message in a chat.
type Message struct {
	ID                  int64           `json:"id" db:"id"`
	ChatID              uuid.UUID       `json:"chat_id" db:"chat_id"`
	SenderID            *uuid.UUID      `json:"sender_id,omitempty" db:"sender_id"`
	ContentType         ContentType     `json:"content_type" db:"content_type"`
	Content             string          `json:"content" db:"content"`
	MediaURL            *string         `json:"media_url,omitempty" db:"media_url"`
	MediaMetadata       json.RawMessage `json:"media_metadata,omitempty" db:"media_metadata"`
	ReplyToMessageID    *int64          `json:"reply_to_message_id,omitempty" db:"reply_to_message_id"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	ContentSearchVector *string         `json:"-" db:"content_search_vector"`
	EditedAt            *time.Time      `json:"edited_at,omitempty" db:"edited_at"`
	DeletedAt           *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedByUserID     *uuid.UUID      `json:"deleted_by_user_id,omitempty" db:"deleted_by_user_id"`
}

// MessageStatusType defines the status of a message for a user.
type MessageStatusType string

const (
	DeliveredStatus MessageStatusType = "delivered"
	ReadStatus      MessageStatusType = "read"
)

// MessageStatus represents the status of a message for a specific user.
type MessageStatus struct {
	MessageID int64             `json:"message_id" db:"message_id"`
	UserID    uuid.UUID         `json:"user_id" db:"user_id"`
	Status    MessageStatusType `json:"status" db:"status"`
	Timestamp time.Time         `json:"timestamp" db:"timestamp"`
}

// MessageReaction represents a reaction to a message.
type MessageReaction struct {
	MessageID     int64     `json:"message_id" db:"message_id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	ReactionEmoji string    `json:"reaction_emoji" db:"reaction_emoji"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
