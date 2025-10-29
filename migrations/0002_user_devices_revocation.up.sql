-- Add revocation support to user_devices
ALTER TABLE user_devices
    ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMPTZ;

-- Helpful indexes for lookups
CREATE INDEX IF NOT EXISTS user_devices_user_id_idx ON user_devices (user_id);
CREATE INDEX IF NOT EXISTS user_devices_revoked_at_idx ON user_devices (revoked_at);
-- Optional: fast lookup by hash
CREATE INDEX IF NOT EXISTS user_devices_refresh_hash_idx ON user_devices (refresh_token_hash);

-- Ensure one record per (user, refresh hash)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'user_devices_user_hash_unique'
    ) THEN
        ALTER TABLE user_devices
        ADD CONSTRAINT user_devices_user_hash_unique UNIQUE (user_id, refresh_token_hash);
    END IF;
END $$;
