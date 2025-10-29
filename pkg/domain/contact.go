package domain

import "github.com/google/uuid"

// Contact represents a contact in a user's address book.
type Contact struct {
	UserID              uuid.UUID  `json:"user_id" db:"user_id"`
	ContactPhoneNumber  string     `json:"contact_phone_number" db:"contact_phone_number"`
	ContactUserID       *uuid.UUID `json:"contact_user_id,omitempty" db:"contact_user_id"`
	DisplayNameOverride string     `json:"display_name_override" db:"display_name_override"`
}
