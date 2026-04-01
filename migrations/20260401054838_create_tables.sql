-- +goose Up
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    created_by UUID NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE group_role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS user_groups (
    user_id UUID NOT NULL REFERENCES users (id),
    group_id UUID NOT NULL REFERENCES groups (id),
    role group_role NOT NULL DEFAULT 'user',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY(user_id, group_id)
);

CREATE TYPE message_type AS ENUM ('group', 'private');

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type message_type NOT NULL,

    sender_id UUID NOT NULL REFERENCES users(id),

    user_id UUID REFERENCES users(id),
    group_id UUID REFERENCES groups(id),

    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (
        (type = 'private' AND user_id IS NOT NULL AND group_id IS NULL) OR
        (type = 'group' AND group_id IS NOT NULL AND user_id IS NULL)
    )
);

-- +goose Down
DROP TABLE IF EXISTS user_groups;

DROP TABLE IF EXISTS messages;

DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS group_role;
DROP TYPE IF EXISTS message_type;

DROP EXTENSION IF EXISTS "citext";
