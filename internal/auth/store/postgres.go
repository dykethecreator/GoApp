package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dykethecreator/GoApp/internal/auth/repository"
	"github.com/dykethecreator/GoApp/pkg/domain"
	"github.com/google/uuid"
)

// UserStore implements the UserRepository interface for PostgreSQL.
type UserStore struct {
	db *sql.DB
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *sql.DB) repository.UserRepository {
	return &UserStore{db: db}
}

// FindByPhoneNumber finds a user by their phone number.
func (s *UserStore) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	query := `SELECT id, phone_number, display_name, profile_picture_url, about_text, last_seen_at, created_at, updated_at FROM users WHERE phone_number = $1`

	user := &domain.User{}
	err := s.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.DisplayName,
		&user.ProfilePictureURL,
		&user.AboutText,
		&user.LastSeenAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found is not an application error, so we return nil, nil.
			// The service layer will handle the logic for creating a new user.
			return nil, nil
		}
		return nil, err // A real database error occurred.
	}

	return user, nil
}

// CreateUser creates a new user in the database.
func (s *UserStore) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (id, phone_number, display_name, profile_picture_url, about_text, last_seen_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	// Set default values for a new user
	user.ID = uuid.New()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.LastSeenAt = now
	// DisplayName can be set to a default or left empty
	if user.DisplayName == "" {
		user.DisplayName = "New User" // Or some other default
	}

	err := s.db.QueryRowContext(ctx, query,
		user.ID,
		user.PhoneNumber,
		user.DisplayName,
		user.ProfilePictureURL,
		user.AboutText,
		user.LastSeenAt,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID finds a user by their ID.
func (s *UserStore) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `SELECT id, phone_number, display_name, profile_picture_url, about_text, last_seen_at, created_at, updated_at FROM users WHERE id = $1`

	user := &domain.User{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.DisplayName,
		&user.ProfilePictureURL,
		&user.AboutText,
		&user.LastSeenAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err // Database error
	}

	return user, nil
}
