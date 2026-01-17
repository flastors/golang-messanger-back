package restapi

import "github.com/go-chi/chi/v5"

type Handler interface {
	RegisterRoutes() chi.Router
}
