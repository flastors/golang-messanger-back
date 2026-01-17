package v1

import (
	"GolangMessanger/internal/controllers/restapi/ports"
	"GolangMessanger/internal/domain/chat"
	"encoding/json"
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
		r.Get("/{id}", http.HandlerFunc(h.GetChat))
		r.Post("/", http.HandlerFunc(h.CreateChat))
		r.Delete("/{id}", http.HandlerFunc(h.DeleteChat))
		r.Post("/{id}/messages", http.HandlerFunc(h.CreateMessage))
	})
}

type CreateChatRequest struct {
	Title string `json:"title"`
}

type CreateChatResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var createChatReq CreateChatRequest
	if err := decoder.Decode(&createChatReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	dto := &chat.NewChatDTO{
		Title: createChatReq.Title,
	}
	c, usecaseErr := h.usecase.CreateChat(dto)
	if usecaseErr != nil {
		if usecaseErr == chat.ErrValidationFailed {
			http.Error(w, "Validation failed", http.StatusBadRequest)
			return
		}
		log.Println(usecaseErr.Error())
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}
	res := &CreateChatResponse{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to marshal chat", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
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

func (h *ChatHandler) GetChat(w http.ResponseWriter, req *http.Request) {
	var limit int
	strChatID := chi.URLParam(req, "id")
	chatID, err := strconv.ParseInt(strChatID, 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	strLimit := req.URL.Query().Get("limit")
	if strLimit != "" {
		var limitParseErr error
		limit, limitParseErr = strconv.Atoi(strLimit)
		if limitParseErr != nil {
			limit = 20
		}
	}
	c, msgs, usecaseErr := h.usecase.GetChat(chatID, limit)
	if usecaseErr != nil {
		if usecaseErr == chat.ErrChatNotFound {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		log.Println(usecaseErr.Error())
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
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, req *http.Request) {
	strChatID := chi.URLParam(req, "id")
	chatID, err := strconv.ParseInt(strChatID, 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	if err := h.usecase.DeleteChat(chatID); err != nil {
		log.Println(err.Error())
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

func (h *ChatHandler) CreateMessage(w http.ResponseWriter, req *http.Request) {
	strChatID := chi.URLParam(req, "id")
	chatID, err := strconv.ParseInt(strChatID, 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(req.Body)
	var createMsgReq CreateMessageRequest
	if err := decoder.Decode(&createMsgReq); err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	dto := &chat.NewMessageDTO{
		ChatID: chatID,
		Text:   createMsgReq.Text,
	}
	m, usecaseErr := h.usecase.CreateMessage(dto)
	if usecaseErr != nil {
		if usecaseErr == chat.ErrValidationFailed {
			log.Println(usecaseErr.Error())
			http.Error(w, "Validation failed", http.StatusBadRequest)
			return
		}
		if usecaseErr == chat.ErrChatNotFound {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		log.Println(usecaseErr.Error())
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}
	res := &CreateMessageResponse{
		ID:        m.ID,
		ChatID:    m.ChatID,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to marshal message", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
