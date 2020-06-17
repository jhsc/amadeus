package api

import (
	"net/http"

	"gitlab.com/jhsc/amadeus/docker"
)

func (h *Handler) handleNewDeployment(w http.ResponseWriter, r *http.Request) {
	req := docker.DeployerPayload{}

	err := h.parseRequest(r, &req)
	if err != nil {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Invalid request body")
		return
	}

	if req.ComposeFile == "" {
		h.renderError(w, http.StatusBadRequest, "BadRequest", "Invalid Compose File")
		return
	}

	err = h.DockerService.DeployCompose(req)
	if err != nil {
		h.logError("payload to deploy: %+v\nError: %s", req, err)
		h.renderError(w, http.StatusInternalServerError, "ServerError", "Server error")
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Succesfully deployed",
	}

	h.render(w, http.StatusCreated, response)
}
