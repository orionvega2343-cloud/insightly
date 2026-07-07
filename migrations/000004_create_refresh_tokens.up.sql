CREATE TABLE IF NOT EXISTS refresh_tokens(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    expires_at TIMESTAMP NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);