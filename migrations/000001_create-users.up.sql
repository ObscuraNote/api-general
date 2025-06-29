CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    user_address CHAR(64) NOT NULL UNIQUE,
    password CHAR(64) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_address_format CHECK (
        user_address ~ '^[0-9a-f]{64}$'
    ),
    CONSTRAINT password_hash_format CHECK (password ~ '^[0-9a-f]{64}$')
);