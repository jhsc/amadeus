package api

import (
	"log"

	"github.com/go-chi/chi"
	"gitlab.com/jhsc/amadeus/docker"
)

// Config is an API handler configuration.
type Config struct {
	Logger        *log.Logger
	DockerService *docker.Service
}

// Handler handles API requests.
type Handler struct {
	*Config
	router chi.Router
}

// New creates a new handler based on the given config.
func New(config *Config) *Handler {
	h := &Handler{Config: config}

	h.router = chi.NewRouter()

	return h
}
