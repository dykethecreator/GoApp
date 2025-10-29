package repository

import (
	"context"

	"github.com/dykethecreator/GoApp/pkg/domain"
)

// UserRepository defines the interface for user data operations.
// This acts as a contract for the data layer (store).
type UserRepository interface {
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByID(ctx context.Context, userID string) (*domain.User, error)
}
