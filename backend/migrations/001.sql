CREATE TABLE IF NOT EXISTS users (
id BIGSERIAL PRIMARY KEY,
name TEXT NOT NULL,
email TEXT UNIQUE NOT NULL,
hashed_password BYTEA NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_movies (
id BIGSERIAL PRIMARY KEY,
user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
movie_id BIGINT NOT NULL,
status TEXT NOT NULL CHECK (status IN ('Plan to Watch', 'Watching', 'Completed', 'Dropped')),
user_rating INTEGER CHECK (user_rating BETWEEN 1 AND 10),
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
UNIQUE(user_id, movie_id)
);

CREATE INDEX IF NOT EXISTS idx_user_movies_user_id ON user_movies (user_id); 