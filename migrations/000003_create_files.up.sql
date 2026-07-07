CREATE TABLE IF NOT EXISTS queries(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    file_id INT NOT NULL REFERENCES files(id),
    created_at TIMESTAMP DEFAULT NOW(),
    question TEXT NOT NULL,
    answer TEXT NOT NULL
);