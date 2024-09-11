CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    firstname   VARCHAR(50) NOT NULL,
    lastname    VARCHAR(50) NOT NULL,
    username    VARCHAR(50) NOT NULL,
    email       VARCHAR(100) NOT NULL,
    password    TEXT NOT NULL,
    phone       VARCHAR(15) NOT NULL,
    activated   BOOLEAN NOT NULL,
    role        VARCHAR(20) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ DEFAULT now(),

    UNIQUE (username, email)
);

create unique index if not exists idx_user_username on users (username);
create unique index if not exists idx_user_email on users (email);
