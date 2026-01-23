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

	_ "GolangMessanger/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
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
	middlewares.RegisterMetrics()
	r := chi.NewRouter()
	// CORS должен идти первым, чтобы preflight не доходил до метрик/логгера
	r.Use(middlewares.CORS)
	r.Use(middlewares.Recoverer)
	r.Use(middlewares.RequestLogger)
	r.Use(middlewares.NewMetricsMiddleware())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.Get("/metrics", middlewares.MetricsHandler().ServeHTTP)
	// Swagger UI route (served by swag-generated docs package + http-swagger)
	r.Handle("/swagger/*", httpSwagger.Handler())
	// Swagger UI available at /swagger/index.html (and /swagger/*)
	r.Route("/api", func(r chi.Router) {
		v1Router := v1.NewRouterV1(usecases.ChatUsecase)
		v1Router.RegisterRoutes(r)
	})
	return r
}
