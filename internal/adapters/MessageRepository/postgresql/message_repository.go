package postgresql

import (
	"GolangMessanger/internal/domain/chat"
	"time"

	"gorm.io/gorm"
)

// messageModel — GORM-модель для таблицы messages
// В ней используются типы, совместимые с доменной моделью chat.Message
type messageModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID    int64     `gorm:"index;not null" json:"chat_id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (messageModel) TableName() string {
	return "messages"
}

// messageRepository реализует domain/chat.MessageRepository
type messageRepository struct {
	db *gorm.DB
}

// NewRepository создаёт новый репозиторий сообщений
func NewRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{db: db}
}

// Create сохраняет сообщение и заполняет его ID/CreatedAt
func (r *messageRepository) Create(m *chat.Message) error {
	model := &messageModel{
		ChatID: m.ChatID,
		Text:   m.Text,
	}
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	if err := r.db.First(model, model.ID).Error; err != nil {
		return err
	}
	m.ID = model.ID
	m.CreatedAt = model.CreatedAt
	return nil
}

// FindByChatID возвращает сообщения для чата, лимитируется limit (если limit<=0 — без лимита)
func (r *messageRepository) FindByChatID(chatID int64, limit int) ([]*chat.Message, error) {
	var models []messageModel
	q := r.db.Where("chat_id = ?", chatID).Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	messages := make([]*chat.Message, 0, len(models))
	for _, mm := range models {
		m := &chat.Message{
			ID:        mm.ID,
			ChatID:    mm.ChatID,
			Text:      mm.Text,
			CreatedAt: mm.CreatedAt,
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// DeleteByChatID удаляет все сообщения чата
func (r *messageRepository) DeleteByChatID(chatID int64) error {
	return r.db.Where("chat_id = ?", chatID).Delete(&messageModel{}).Error
}
