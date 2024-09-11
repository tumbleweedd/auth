CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS outbox
(
    id         uuid        NOT NULL DEFAULT uuid_generate_v4(),
    event_type VARCHAR(50) NOT NULL,
    payload    JSONB       NOT NULL,
    created_at timestamp            default now(),

    PRIMARY KEY (id)
)