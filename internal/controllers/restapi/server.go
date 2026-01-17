package restapi

import (
	"GolangMessanger/internal/config"
	"GolangMessanger/internal/controllers/restapi/middlewares"
	v1 "GolangMessanger/internal/controllers/restapi/v1"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	server *http.Server
	cfg    config.RestAPIConfig
}

func NewServer(usecases Usecases, cfg config.RestAPIConfig) *Server {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      Router(usecases),
	}
	return &Server{
		server: server,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	log.Println("Starting server on ", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping server...")
	return s.server.Shutdown(ctx)
}

func Router(usecases Usecases) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.Recoverer)
	r.Use(middlewares.RequestLogger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Route("/api", func(r chi.Router) {
		v1Router := v1.NewRouterV1(usecases.ChatUsecase)
		v1Router.RegisterRoutes(r)
	})
	return r
}
