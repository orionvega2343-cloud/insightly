CREATE  TABLE IF NOT EXISTS files(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    filename TEXT NOT NULL,
    filepath TEXT NOT NULL
)