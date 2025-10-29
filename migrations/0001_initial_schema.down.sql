-- This is the down migration for the initial schema.
-- It drops all the tables created in the up migration.

DROP TABLE IF EXISTS poll_votes;
DROP TABLE IF EXISTS poll_options;
DROP TABLE IF EXISTS polls;
DROP TABLE IF EXISTS user_devices;
DROP TABLE IF EXISTS call_logs;
DROP TABLE IF EXISTS message_reactions;
DROP TABLE IF EXISTS blocked_users;
DROP TABLE IF EXISTS status_views;
DROP TABLE IF EXISTS status_updates;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS message_statuses;
ALTER TABLE chats DROP CONSTRAINT IF EXISTS fk_pinned_message;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chat_members;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;
