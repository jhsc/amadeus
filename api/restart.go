package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) handleRestartContainer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(h.urlParam(r, "id"), 10, 64)
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Invalid compose id")
		return
	}
	srv := r.URL.Query().Get("service")

	msg := fmt.Sprintf("Succesfully restarted ComposeId: %d, service: %s", id, srv)
	response := struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}

	h.render(w, http.StatusCreated, response)
}
