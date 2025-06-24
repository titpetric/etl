-- petstore.sql
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS pets;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS orders;

CREATE TABLE pets (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    age INTEGER,
    status TEXT NOT NULL -- available, sold, pending
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT
);

CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    quantity INTEGER DEFAULT 1,
    order_date TEXT NOT NULL,
    status TEXT NOT NULL, -- placed, approved, delivered
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

-- Seed data
INSERT INTO pets (id, name, type, age, status) VALUES
(1, 'Fluffy', 'Cat', 3, 'available'),
(2, 'Spot', 'Dog', 5, 'sold'),
(3, 'Goldie', 'Fish', 1, 'available');

INSERT INTO users (id, username, email) VALUES
(1, 'alice', 'alice@example.com'),
(2, 'bob', 'bob@example.com');

INSERT INTO orders (id, user_id, pet_id, quantity, order_date, status) VALUES
(1, 1, 2, 1, '2024-06-01T10:00:00Z', 'delivered'),
(2, 2, 1, 2, '2024-06-10T14:30:00Z', 'placed');
