CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_secrets
(
    user_id            UUID PRIMARY KEY,

    password_hash      TEXT NOT NULL,
    password_salt      TEXT NOT NULL,

    refresh_token_hash TEXT NOT NULL,
    refresh_token_salt TEXT NOT NULL,

    created_at         TIMESTAMPTZ DEFAULT now(),
    updated_at         TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
)