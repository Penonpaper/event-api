

CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(120) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(120) NOT NULL DEFAULT 'client',
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE events(
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description text,
    total_seats INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);