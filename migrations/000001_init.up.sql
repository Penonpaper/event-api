

CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(120) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(120) NOT NULL DEFAULT 'client',
    created_at TIMESTAMP DEFAULT NOW(),
    nickname VARCHAR(100) UNIQUE
);
CREATE TABLE events(
    id UUID PRIMARY KEY,
    organizer_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description text,
    location VARCHAR(160) NOT NULL,
    status VARCHAR(50) NOT NULL,
    start_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    total_seats INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK(ends_at > start_at)

);