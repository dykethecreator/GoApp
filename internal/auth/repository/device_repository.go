package repository

import (
	"context"

	"github.com/dykethecreator/GoApp/pkg/domain"
)

// DeviceRepository defines operations for user device sessions (refresh tokens).
type DeviceRepository interface {
	UpsertDevice(ctx context.Context, dev *domain.UserDevice) error
	FindActiveByUserAndHash(ctx context.Context, userID string, hash string) (*domain.UserDevice, error)
	RevokeByID(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
