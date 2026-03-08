package v1

import (
	"GolangMessanger/internal/controllers/restapi/ports"
	"GolangMessanger/internal/domain/chat"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ChatHandler struct {
	usecase ports.ChatUsecase
}

func NewChatHandler(usecase ports.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
	}
}

func (h *ChatHandler) RegisterRoutes(r chi.Router) {
	r.Route("/chats", func(r chi.Router) {
		r.Get("/{id}", h.GetChat)
		r.Post("/", h.CreateChat)
		r.Delete("/{id}", h.DeleteChat)
		r.Post("/{id}/messages", h.CreateMessage)
	})
}

// writeJSON writes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("writeJSON encode error:", err)
	}
}

type CreateChatRequest struct {
	Title string `json:"title"`
}

type CreateChatResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateChat godoc
// @Summary Create chat
// @Description Create a new chat
// @Tags chats
// @Accept json
// @Produce json
// @Param payload body CreateChatRequest true "Create chat request"
// @Success 200 {object} CreateChatResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /chats [post]
func (h *ChatHandler) CreateChat(w http.ResponseWriter, req *http.Request) {
	var createChatReq CreateChatRequest
	if err := json.NewDecoder(req.Body).Decode(&createChatReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	c, err := h.usecase.CreateChat(&chat.NewChatDTO{Title: createChatReq.Title})
	if err != nil {
		if errors.Is(err, chat.ErrValidationFailed) {
			http.Error(w, "Validation failed", http.StatusBadRequest)
			return
		}
		log.Println(err)
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, &CreateChatResponse{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
	})
}

type GetChatResponse struct {
	ID        int64                     `json:"id"`
	Title     string                    `json:"title"`
	CreatedAt time.Time                 `json:"created_at"`
	Messages  []*GetChatMessageResponse `json:"messages"`
}

type GetChatMessageResponse struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// GetChat godoc
// @Summary Get chat with messages
// @Tags chats
// @Param id path int true "Chat ID"
// @Param limit query int false "Limit messages"
// @Produce json
// @Success 200 {object} GetChatResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /chats/{id} [get]
func (h *ChatHandler) GetChat(w http.ResponseWriter, req *http.Request) {
	chatID, err := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	limit := 0 // 0 means "no explicit limit" — the usecase enforces MaxMessageLimit
	if s := req.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			limit = n
		}
	}
	c, msgs, err := h.usecase.GetChat(chatID, limit)
	if err != nil {
		if errors.Is(err, chat.ErrChatNotFound) {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "Failed to get chat", http.StatusInternalServerError)
		return
	}
	res := &GetChatResponse{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
		Messages:  make([]*GetChatMessageResponse, 0, len(msgs)),
	}
	for _, msg := range msgs {
		res.Messages = append(res.Messages, &GetChatMessageResponse{
			ID:        msg.ID,
			Text:      msg.Text,
			CreatedAt: msg.CreatedAt,
		})
	}
	writeJSON(w, http.StatusOK, res)
}

// DeleteChat godoc
// @Summary Delete chat
// @Tags chats
// @Param id path int true "Chat ID"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /chats/{id} [delete]
func (h *ChatHandler) DeleteChat(w http.ResponseWriter, req *http.Request) {
	chatID, err := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	if err := h.usecase.DeleteChat(chatID); err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete chat", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type CreateMessageRequest struct {
	Text string `json:"text"`
}

type CreateMessageResponse struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateMessage godoc
// @Summary Create message in chat
// @Tags messages
// @Accept json
// @Param id path int true "Chat ID"
// @Param payload body CreateMessageRequest true "Create message request"
// @Produce json
// @Success 200 {object} CreateMessageResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /chats/{id}/messages [post]
func (h *ChatHandler) CreateMessage(w http.ResponseWriter, req *http.Request) {
	chatID, err := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	var createMsgReq CreateMessageRequest
	if err := json.NewDecoder(req.Body).Decode(&createMsgReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	m, err := h.usecase.CreateMessage(&chat.NewMessageDTO{
		ChatID: chatID,
		Text:   createMsgReq.Text,
	})
	if err != nil {
		if errors.Is(err, chat.ErrValidationFailed) {
			http.Error(w, "Validation failed", http.StatusBadRequest)
			return
		}
		if errors.Is(err, chat.ErrChatNotFound) {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, &CreateMessageResponse{
		ID:        m.ID,
		ChatID:    m.ChatID,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
	})
}
