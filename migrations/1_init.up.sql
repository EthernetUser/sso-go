CREATE TABLE
    IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        is_admin BOOLEAN DEFAULT FALSE
    );

CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);

CREATE TABLE
    IF NOT EXISTS apps (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        secret TEXT NOT NULL
    );