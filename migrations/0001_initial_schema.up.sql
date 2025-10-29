-- This is the initial schema for the database.
-- It creates all the tables based on the provided schema.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    profile_picture_url TEXT,
    about_text VARCHAR(250),
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE chats (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL, -- 'one_to_one' or 'group'
    group_name VARCHAR(100),
    group_icon_url TEXT,
    group_description TEXT,
    created_by_user_id uuid REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_message_at TIMESTAMPTZ,
    pinned_message_id BIGINT
);

CREATE TABLE chat_members (
    chat_id uuid NOT NULL REFERENCES chats(id),
    user_id uuid NOT NULL REFERENCES users(id),
    role VARCHAR(20) NOT NULL, -- 'admin' or 'member'
    membership_status VARCHAR(20) NOT NULL, -- 'active', 'left', 'kicked'
    is_muted BOOLEAN DEFAULT FALSE,
    is_archived BOOLEAN DEFAULT FALSE,
    unread_count INTEGER DEFAULT 0,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id uuid NOT NULL REFERENCES chats(id),
    sender_id uuid REFERENCES users(id),
    content_type VARCHAR(30) NOT NULL,
    content TEXT,
    media_url TEXT,
    media_metadata JSONB,
    reply_to_message_id BIGINT REFERENCES messages(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    content_search_vector TSVECTOR,
    edited_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    deleted_by_user_id uuid REFERENCES users(id)
);

ALTER TABLE chats ADD CONSTRAINT fk_pinned_message FOREIGN KEY (pinned_message_id) REFERENCES messages(id);
CREATE INDEX ON messages (chat_id, created_at DESC);
CREATE INDEX ON messages USING GIN(content_search_vector);


CREATE TABLE message_statuses (
    message_id BIGINT NOT NULL REFERENCES messages(id),
    user_id uuid NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL, -- 'delivered' or 'read'
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (message_id, user_id)
);

CREATE TABLE contacts (
    user_id uuid NOT NULL REFERENCES users(id),
    contact_phone_number VARCHAR(20) NOT NULL,
    contact_user_id uuid REFERENCES users(id),
    display_name_override VARCHAR(100),
    PRIMARY KEY (user_id, contact_phone_number)
);

CREATE TABLE status_updates (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL REFERENCES users(id),
    content_type VARCHAR(20) NOT NULL, -- 'text', 'image', 'video'
    content_or_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX ON status_updates (user_id, created_at);


CREATE TABLE status_views (
    status_id uuid NOT NULL REFERENCES status_updates(id),
    user_id uuid NOT NULL REFERENCES users(id),
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (status_id, user_id)
);

CREATE TABLE blocked_users (
    blocker_user_id uuid NOT NULL REFERENCES users(id),
    blocked_user_id uuid NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (blocker_user_id, blocked_user_id)
);

CREATE TABLE message_reactions (
    message_id BIGINT NOT NULL REFERENCES messages(id),
    user_id uuid NOT NULL REFERENCES users(id),
    reaction_emoji VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (message_id, user_id)
);

CREATE TABLE call_logs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id uuid NOT NULL REFERENCES chats(id),
    caller_user_id uuid NOT NULL REFERENCES users(id),
    call_type VARCHAR(10) NOT NULL, -- 'audio' or 'video'
    status VARCHAR(20) NOT NULL, -- 'completed', 'missed', 'rejected'
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ
    -- Postgres does not support inline INDEX; create it separately below
);

-- Create index for faster lookups by chat_id
CREATE INDEX IF NOT EXISTS call_logs_chat_id_idx ON call_logs (chat_id);

CREATE TABLE user_devices (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL REFERENCES users(id),
    refresh_token_hash VARCHAR(255) NOT NULL,
    device_name VARCHAR(100),
    device_type VARCHAR(20), -- 'mobile', 'web', 'desktop'
    push_notification_token TEXT,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE polls (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id uuid NOT NULL REFERENCES chats(id),
    created_by_user_id uuid NOT NULL REFERENCES users(id),
    question_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE poll_options (
    id BIGSERIAL PRIMARY KEY,
    poll_id uuid NOT NULL REFERENCES polls(id),
    option_text VARCHAR(255) NOT NULL
);

CREATE TABLE poll_votes (
    poll_id uuid NOT NULL REFERENCES polls(id),
    option_id BIGINT NOT NULL REFERENCES poll_options(id),
    user_id uuid NOT NULL REFERENCES users(id),
    PRIMARY KEY (poll_id, user_id)
);
