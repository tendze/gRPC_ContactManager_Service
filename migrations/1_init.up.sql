CREATE TABLE IF NOT EXISTS contacts(
    id INTEGER PRIMARY KEY,
    creator_email TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL
);