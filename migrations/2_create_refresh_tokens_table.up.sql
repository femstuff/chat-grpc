CREATE TABLE IF NOT EXISTS refresh_tokens (
      user_id BIGINT PRIMARY KEY,
      token TEXT NOT NULL
);
