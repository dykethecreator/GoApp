package domain

import (
	"time"

	"github.com/google/uuid"
)

// ChatType defines the type of chat.
type ChatType string

const (
	OneToOneChat ChatType = "one_to_one"
	GroupChat    ChatType = "group"
)

// Chat represents a conversation, either one-to-one or a group.
type Chat struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Type             ChatType  `json:"type" db:"type"`
	GroupName        *string   `json:"group_name,omitempty" db:"group_name"`
	GroupIconURL     *string   `json:"group_icon_url,omitempty" db:"group_icon_url"`
	GroupDescription *string   `json:"group_description,omitempty" db:"group_description"`
	CreatedByUserID  uuid.UUID `json:"created_by_user_id" db:"created_by_user_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	LastMessageAt    time.Time `json:"last_message_at" db:"last_message_at"`
	PinnedMessageID  *int64    `json:"pinned_message_id,omitempty" db:"pinned_message_id"`
}

// ChatMemberRole defines the role of a member in a chat.
type ChatMemberRole string

const (
	AdminRole  ChatMemberRole = "admin"
	MemberRole ChatMemberRole = "member"
)

// MembershipStatus defines the status of a chat member.
type MembershipStatus string

const (
	ActiveMembership MembershipStatus = "active"
	LeftMembership   MembershipStatus = "left"
	KickedMembership MembershipStatus = "kicked"
)

// ChatMember represents a user's membership in a chat.
type ChatMember struct {
	ChatID           uuid.UUID        `json:"chat_id" db:"chat_id"`
	UserID           uuid.UUID        `json:"user_id" db:"user_id"`
	Role             ChatMemberRole   `json:"role" db:"role"`
	MembershipStatus MembershipStatus `json:"membership_status" db:"membership_status"`
	IsMuted          bool             `json:"is_muted" db:"is_muted"`
	IsArchived       bool             `json:"is_archived" db:"is_archived"`
	UnreadCount      int              `json:"unread_count" db:"unread_count"`
	JoinedAt         time.Time        `json:"joined_at" db:"joined_at"`
}
