package api

import (
	"net/http"
)

// HandleHealthz renders 200 status if server is running
func (h *Handler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	id, err := h.Store.Projects().New("Test Project")
	if err != nil {
		h.logError("create new project: %s", err)
		h.renderError(w, http.StatusInternalServerError, "ServerError", "Server error")
		return
	}
	h.Logger.Printf("Created project with id: %d\n", id)
	response := struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	}

	h.render(w, http.StatusOK, response)
}
