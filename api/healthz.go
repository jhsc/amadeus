package api

import (
	"net/http"
)

func (h *Handler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	}

	h.render(w, http.StatusOK, response)
}
