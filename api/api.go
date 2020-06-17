package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"

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
	Token  string
}

// New creates a new handler based on the given config.
func New(config *Config, token string) *Handler {
	h := &Handler{Config: config, Token: token}

	h.router = chi.NewRouter()
	h.router.Use(h.restrictAccess)

	// h.router.Get("/healthz", h.handleHealthz)

	// TODO: create deployment handle
	// h.router.Post("/deploy", h.handleNewDeployment)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) urlParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (h *Handler) parseRequest(r *http.Request, data interface{}) error {
	const maxRequestLen = 16 * 1024 * 1024
	lr := io.LimitReader(r.Body, maxRequestLen)
	return json.NewDecoder(lr).Decode(data)
}

func (h *Handler) render(w http.ResponseWriter, status int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logError("marshal json: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func (h *Handler) renderError(w http.ResponseWriter, status int, code, message string) {
	response := struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{}
	response.Error.Code = code
	response.Error.Message = message
	h.render(w, status, response)
}

func (h *Handler) logError(format string, a ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	callerNameSplit := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	funcName := callerNameSplit[len(callerNameSplit)-1]
	h.Logger.Printf("ERROR: %s: %s", funcName, fmt.Sprintf(format, a...))
}

// restrictAccess middleware to allow request by valid token only
func (h *Handler) restrictAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == h.Token {
			next.ServeHTTP(w, r)
			return
		}

		h.renderError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), "Not authorized")
		return
	})
}
