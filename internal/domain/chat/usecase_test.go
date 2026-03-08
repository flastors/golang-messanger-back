package chat

import (
	"errors"
	"testing"

	"GolangMessanger/internal/config"

	"github.com/stretchr/testify/require"
)

// Тесты при необходимости переопределяют поведение методов репозиториев
type mockChatRepo struct {
	createFn   func(*Chat) error
	findByIDFn func(int64) (*Chat, error)
	deleteFn   func(int64) error
}

func (m *mockChatRepo) Create(c *Chat) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(c)
}

func (m *mockChatRepo) FindByID(id int64) (*Chat, error) {
	if m.findByIDFn == nil {
		return nil, nil
	}
	return m.findByIDFn(id)
}

func (m *mockChatRepo) Delete(id int64) error {
	if m.deleteFn == nil {
		return nil
	}
	return m.deleteFn(id)
}

type mockMessageRepo struct {
	createFn       func(*Message) error
	findByChatIDFn func(int64, int) ([]*Message, error)
}

func (m *mockMessageRepo) Create(msg *Message) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(msg)
}

func (m *mockMessageRepo) FindByChatID(chatID int64, limit int) ([]*Message, error) {
	if m.findByChatIDFn == nil {
		return nil, nil
	}
	return m.findByChatIDFn(chatID, limit)
}

func makeUsecaseWithCfg() *chatUsecase {
	// тесты при необходимости переопределяют
	cfg := config.ChatServiceConfig{
		MinTitleSize:    1,
		MaxTitleSize:    100,
		MaxMessageLimit: 50,
		MaxMessageSize:  1000,
	}
	return &chatUsecase{
		chatRepo:    &mockChatRepo{},
		messageRepo: &mockMessageRepo{},
		cfg:         cfg,
	}
}

func TestCreateChat_ValidationAndSuccess(t *testing.T) {
	cfg := config.ChatServiceConfig{MinTitleSize: 1, MaxTitleSize: 5}
	t.Run("nil dto", func(t *testing.T) {
		u := makeUsecaseWithCfg()
		u.cfg = cfg
		_, err := u.CreateChat(nil)
		require.Error(t, err)
	})

	t.Run("too short title", func(t *testing.T) {
		u := makeUsecaseWithCfg()
		u.cfg = cfg
		_, err := u.CreateChat(&NewChatDTO{Title: "a"})
		require.Error(t, err)
	})

	t.Run("too long title", func(t *testing.T) {
		u := makeUsecaseWithCfg()
		u.cfg = cfg
		_, err := u.CreateChat(&NewChatDTO{Title: "abcdef"})
		require.Error(t, err)
	})

	t.Run("success and trim", func(t *testing.T) {
		mock := &mockChatRepo{
			createFn: func(c *Chat) error {
				c.ID = 52
				return nil
			},
		}
		u := makeUsecaseWithCfg()
		u.chatRepo = mock
		u.cfg = cfg

		dto := NewChatDTO{Title: "  hi  "}
		chat, err := u.CreateChat(&dto)
		require.NoError(t, err)
		require.NotNil(t, chat)
		require.Equal(t, int64(52), chat.ID)
		require.Equal(t, "hi", chat.Title)
	})
}

func TestGetChat_Behavior(t *testing.T) {
	t.Run("chat not found", func(t *testing.T) {
		mockC := &mockChatRepo{
			findByIDFn: func(id int64) (*Chat, error) {
				return nil, errors.New("not found")
			},
		}
		u := makeUsecaseWithCfg()
		u.chatRepo = mockC

		_, _, err := u.GetChat(1, 10)
		require.Error(t, err)
	})

	t.Run("limit capped and messages returned", func(t *testing.T) {
		cfg := config.ChatServiceConfig{MaxMessageLimit: 3}
		var receivedLimit int
		mockC := &mockChatRepo{
			findByIDFn: func(id int64) (*Chat, error) {
				return &Chat{ID: id, Title: "t"}, nil
			},
		}
		mockM := &mockMessageRepo{
			findByChatIDFn: func(chatID int64, limit int) ([]*Message, error) {
				receivedLimit = limit
				return []*Message{{ID: 1, ChatID: chatID, Text: "m1"}}, nil
			},
		}
		u := makeUsecaseWithCfg()
		u.cfg = cfg
		u.chatRepo = mockC
		u.messageRepo = mockM

		chat, msgs, err := u.GetChat(7, 10) // 10 > MaxMessageLimit (3)
		require.NoError(t, err)
		require.Equal(t, int64(7), chat.ID)
		require.Equal(t, 1, len(msgs))
		require.Equal(t, cfg.MaxMessageLimit, receivedLimit)
	})
}

func TestDeleteChat(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		mockC := &mockChatRepo{
			deleteFn: func(id int64) error { return errors.New("boom") },
		}
		u := makeUsecaseWithCfg()
		u.chatRepo = mockC

		err := u.DeleteChat(5)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		called := false
		mockC := &mockChatRepo{
			deleteFn: func(id int64) error {
				called = true
				return nil
			},
		}
		u := makeUsecaseWithCfg()
		u.chatRepo = mockC

		err := u.DeleteChat(5)
		require.NoError(t, err)
		require.True(t, called)
	})
}

func TestCreateMessage(t *testing.T) {
	t.Run("validation nil dto", func(t *testing.T) {
		u := makeUsecaseWithCfg()
		_, err := u.CreateMessage(nil)
		require.Error(t, err)
	})

	t.Run("validation text length", func(t *testing.T) {
		cfg := config.ChatServiceConfig{MaxMessageSize: 2}
		u := makeUsecaseWithCfg()
		u.cfg = cfg
		_, err := u.CreateMessage(&NewMessageDTO{ChatID: 1, Text: "too long"})
		require.Error(t, err)
	})

	t.Run("chat not found", func(t *testing.T) {
		mockC := &mockChatRepo{
			findByIDFn: func(id int64) (*Chat, error) { return nil, errors.New("not found") },
		}
		u := makeUsecaseWithCfg()
		u.chatRepo = mockC

		_, err := u.CreateMessage(&NewMessageDTO{ChatID: 1, Text: "ok"})
		require.Error(t, err)
	})

	t.Run("success create message", func(t *testing.T) {
		mockC := &mockChatRepo{
			findByIDFn: func(id int64) (*Chat, error) { return &Chat{ID: id, Title: "t"}, nil },
		}
		mockM := &mockMessageRepo{
			createFn: func(m *Message) error {
				m.ID = 99
				return nil
			},
		}
		u := makeUsecaseWithCfg()
		u.chatRepo = mockC
		u.messageRepo = mockM
		u.cfg = config.ChatServiceConfig{MaxMessageSize: 1000}

		msg, err := u.CreateMessage(&NewMessageDTO{ChatID: 2, Text: "hello"})
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, int64(99), msg.ID)
		require.Equal(t, int64(2), msg.ChatID)
		require.Equal(t, "hello", msg.Text)
	})
}
