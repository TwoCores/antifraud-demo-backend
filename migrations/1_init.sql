-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    -- use varchar here to use IDs from provided dataset for model training
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT,
    last_name TEXT,
    status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    number TEXT NOT NULL,
    balance DOUBLE PRECISION DEFAULT 0,
    status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS login_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    when_ts TIMESTAMP WITH TIME ZONE NOT NULL,
    phone_model TEXT,
    os TEXT
);

CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    from_card_id UUID REFERENCES cards(id),
    to_card_id UUID REFERENCES cards(id),
    amount DOUBLE PRECISION NOT NULL,
    when_ts TIMESTAMP WITH TIME ZONE NOT NULL,
    fraud_score DOUBLE PRECISION,
    is_blocked BOOLEAN DEFAULT FALSE
);
