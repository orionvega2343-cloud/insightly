CREATE TABLE IF NOT EXISTS queries(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    file_id INT NOT NULL REFERENCES files(id),
    question TEXT NOT NULL,
    answer TEXT NOT NULL
)