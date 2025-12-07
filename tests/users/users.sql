-- users.sql
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS user_auth;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_auth (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10,2) NOT NULL,
    status TEXT DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Seed data
INSERT INTO users (id, name, email) VALUES
(1, 'Alice Johnson', 'alice@example.com'),
(2, 'Bob Smith', 'bob@example.com'),
(3, 'Carol White', 'carol@example.com');

INSERT INTO user_auth (user_id, password_hash) VALUES
(1, '$2a$10$hashed_password_alice'),
(2, '$2a$10$hashed_password_bob'),
(3, '$2a$10$hashed_password_carol');

INSERT INTO orders (id, user_id, order_date, total_amount, status) VALUES
(1, 1, '2024-01-15', 99.99, 'completed'),
(2, 1, '2024-02-20', 149.50, 'completed'),
(3, 2, '2024-03-10', 75.00, 'pending'),
(4, 3, '2024-04-05', 225.75, 'completed');
