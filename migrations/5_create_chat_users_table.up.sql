CREATE TABLE IF NOT EXISTS chat_users (
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id INT NOT NULL,
    PRIMARY KEY (chat_id, user_id)
);