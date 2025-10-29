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

// UserDeviceStore implements DeviceRepository for PostgreSQL.
type UserDeviceStore struct {
	db *sql.DB
}

func NewUserDeviceStore(db *sql.DB) repository.DeviceRepository {
	return &UserDeviceStore{db: db}
}

func (s *UserDeviceStore) UpsertDevice(ctx context.Context, dev *domain.UserDevice) error {
	if dev.ID == uuid.Nil {
		dev.ID = uuid.New()
	}
	if dev.CreatedAt.IsZero() {
		dev.CreatedAt = time.Now()
	}
	q := `
	INSERT INTO user_devices (id, user_id, refresh_token_hash, device_name, device_type, push_notification_token, last_login_at, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT (user_id, refresh_token_hash)
	DO UPDATE SET last_login_at = EXCLUDED.last_login_at
	`
	_, err := s.db.ExecContext(ctx, q,
		dev.ID,
		dev.UserID,
		dev.RefreshTokenHash,
		dev.DeviceName,
		dev.DeviceType,
		dev.PushNotificationToken,
		dev.LastLoginAt,
		dev.CreatedAt,
	)
	return err
}

func (s *UserDeviceStore) FindActiveByUserAndHash(ctx context.Context, userID string, hash string) (*domain.UserDevice, error) {
	q := `SELECT id, user_id, refresh_token_hash, device_name, device_type, push_notification_token, last_login_at, created_at, revoked_at
	FROM user_devices WHERE user_id = $1 AND refresh_token_hash = $2 AND revoked_at IS NULL LIMIT 1`
	row := s.db.QueryRowContext(ctx, q, userID, hash)
	var d domain.UserDevice
	if err := row.Scan(&d.ID, &d.UserID, &d.RefreshTokenHash, &d.DeviceName, &d.DeviceType, &d.PushNotificationToken, &d.LastLoginAt, &d.CreatedAt, &d.RevokedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (s *UserDeviceStore) RevokeByID(ctx context.Context, id string) error {
	q := `UPDATE user_devices SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`
	_, err := s.db.ExecContext(ctx, q, id)
	return err
}

func (s *UserDeviceStore) RevokeAllForUser(ctx context.Context, userID string) error {
	q := `UPDATE user_devices SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := s.db.ExecContext(ctx, q, userID)
	return err
}
