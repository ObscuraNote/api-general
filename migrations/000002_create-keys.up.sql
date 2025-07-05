CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS keys (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    user_address CHAR(64) NOT NULL,
    encrypted_key BYTEA NOT NULL,
    key_iv BYTEA NOT NULL,
    encrypted_data BYTEA NOT NULL,
    data_iv BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_keys_user_id ON keys (user_id);