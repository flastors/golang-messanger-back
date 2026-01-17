-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chats_created_at ON chats (created_at);

CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
                          text TEXT NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_chat_id ON messages (chat_id);
CREATE INDEX idx_messages_created_at ON messages (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_chat_id;
DROP INDEX IF EXISTS idx_chats_created_at;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd
