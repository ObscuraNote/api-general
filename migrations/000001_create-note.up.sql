CREATE TABLE keys (
    id uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4 (),
    user_address VARCHAR(64) NOT NULL,
    key BYTEA NOT NULL,
    encrypted_data BYTEA NOT NULL,
    iv BYTEA NOT NULL
);

CREATE INDEX idx_keys_user_address ON keys (user_address);