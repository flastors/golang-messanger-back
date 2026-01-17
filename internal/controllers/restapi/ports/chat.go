package ports

import "GolangMessanger/internal/domain/chat"

type ChatUsecase interface {
	CreateChat(dto *chat.NewChatDTO) (*chat.Chat, error)
	GetChat(chatID int64, limit int) (*chat.Chat, []*chat.Message, error)
	DeleteChat(chatID int64) error
	CreateMessage(dto *chat.NewMessageDTO) (*chat.Message, error)
}
