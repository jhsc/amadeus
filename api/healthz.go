package api

import (
	"net/http"
)

// HandleHealthz renders 200 status if server is running
func (h *Handler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	}

	h.render(w, http.StatusOK, response)
}
