package chat

type NewChatDTO struct {
	Title string `json:"title"`
}

type NewMessageDTO struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}
