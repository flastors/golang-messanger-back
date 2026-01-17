package chat

import "errors"

var (
	ErrChatNotFound     error = errors.New("chat not found")
	ErrValidationFailed       = errors.New("validation failed")
)

type ChatRepository interface {
	Create(chat *Chat) error
	FindByID(chatID int64) (*Chat, error)
	Delete(chatID int64) error
}

type MessageRepository interface {
	Create(message *Message) error
	FindByChatID(chatID int64, limit int) ([]*Message, error)
}
