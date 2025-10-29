-- Revert revocation changes
ALTER TABLE user_devices
    DROP COLUMN IF EXISTS revoked_at;

DROP INDEX IF EXISTS user_devices_user_id_idx;
DROP INDEX IF EXISTS user_devices_revoked_at_idx;
DROP INDEX IF EXISTS user_devices_refresh_hash_idx;
