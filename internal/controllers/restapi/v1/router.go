package v1

import (
	"GolangMessanger/internal/controllers/restapi/ports"

	"github.com/go-chi/chi/v5"
)

type RouterV1 struct {
	ChatUsecase ports.ChatUsecase
}

func NewRouterV1(chatUsecase ports.ChatUsecase) *RouterV1 {
	return &RouterV1{
		ChatUsecase: chatUsecase,
	}
}

func (r *RouterV1) RegisterRoutes(router chi.Router) {
	router.Route("/v1", func(router chi.Router) {
		chatHandler := NewChatHandler(r.ChatUsecase)
		chatHandler.RegisterRoutes(router)
	})
}
