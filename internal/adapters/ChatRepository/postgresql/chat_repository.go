package postgresql

import (
	"GolangMessanger/internal/domain/chat"
	"time"

	"gorm.io/gorm"
)

// chatModel — GORM-модель для таблицы chats
type chatModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(200);not null" json:"title"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (chatModel) TableName() string {
	return "chats"
}

// chatRepository реализует domain/chat.ChatRepository
type chatRepository struct {
	db *gorm.DB
}

// NewRepository создаёт репозиторий чатов
func NewRepository(db *gorm.DB) *chatRepository {
	return &chatRepository{db: db}
}

// Create сохраняет чат и заполняет ID и CreatedAt
func (r *chatRepository) Create(c *chat.Chat) error {
	model := &chatModel{
		Title: c.Title,
	}
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	if err := r.db.First(model, model.ID).Error; err != nil {
		return err
	}
	c.ID = model.ID
	c.CreatedAt = model.CreatedAt
	return nil
}

// FindByID возвращает чат по ID
func (r *chatRepository) FindByID(chatID int64) (*chat.Chat, error) {
	var model chatModel
	if err := r.db.First(&model, chatID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, chat.ErrChatNotFound
		}
		return nil, err
	}
	return &chat.Chat{
		ID:        model.ID,
		Title:     model.Title,
		CreatedAt: model.CreatedAt,
	}, nil
}

// Delete удаляет чат по ID
func (r *chatRepository) Delete(chatID int64) error {
	if err := r.db.Delete(&chatModel{}, chatID).Error; err != nil {
		return err
	}
	return nil
}
