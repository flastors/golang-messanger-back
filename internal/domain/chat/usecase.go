package chat

import (
	"GolangMessanger/internal/config"
	"strings"
	"unicode/utf8"
)

type chatUsecase struct {
	ChatRepo    ChatRepository
	MessageRepo MessageRepository
	cfg         config.ChatServiceConfig
}

func NewChatUsecase(chatRepo ChatRepository, messageRepo MessageRepository, config config.ChatServiceConfig) *chatUsecase {
	return &chatUsecase{
		ChatRepo:    chatRepo,
		MessageRepo: messageRepo,
		cfg:         config,
	}
}

// Бизнес-логика создания чата
func (u *chatUsecase) CreateChat(dto *NewChatDTO) (*Chat, error) {
	if dto == nil {
		return nil, ErrValidationFailed
	}
	trimTitle := strings.TrimSpace(dto.Title)
	if trimTitle == "" || utf8.RuneCountInString(trimTitle) <= u.cfg.MinTitleSize || utf8.RuneCountInString(trimTitle) >= u.cfg.MaxTitleSize {
		return nil, ErrValidationFailed
	}
	chat := &Chat{
		Title: trimTitle,
	}
	if err := u.ChatRepo.Create(chat); err != nil {
		return nil, err
	}
	return chat, nil
}

// Бизнес-логика получения чата и сообщений
func (u *chatUsecase) GetChat(chatID int64, limit int) (*Chat, []*Message, error) {
	if limit > u.cfg.MaxMessageLimit {
		limit = u.cfg.MaxMessageLimit
	}
	chat, chatErr := u.ChatRepo.FindByID(chatID)
	if chatErr != nil {
		return nil, nil, chatErr
	}
	messages, msgErr := u.MessageRepo.FindByChatID(chatID, limit)
	if msgErr != nil {
		return nil, nil, msgErr
	}
	return chat, messages, nil
}

// Бизнес-логика удаления чата (удаление сообщений реализуется на уровне БД через каскадное удаление)
func (u *chatUsecase) DeleteChat(chatID int64) error {

	if err := u.ChatRepo.Delete(chatID); err != nil {
		return err
	}
	return nil
}

// Бизнес-логика создания сообщения
func (u *chatUsecase) CreateMessage(dto *NewMessageDTO) (*Message, error) {
	if dto == nil {
		return nil, ErrValidationFailed
	}
	if dto.Text == "" || utf8.RuneCountInString(dto.Text) < 1 || utf8.RuneCountInString(dto.Text) > u.cfg.MaxMessageSize {
		return nil, ErrValidationFailed
	}
	_, chatErr := u.ChatRepo.FindByID(dto.ChatID)
	if chatErr != nil {
		return nil, chatErr
	}
	message := &Message{
		ChatID: dto.ChatID,
		Text:   dto.Text,
	}
	if err := u.MessageRepo.Create(message); err != nil {
		return nil, err
	}
	return message, nil
}
