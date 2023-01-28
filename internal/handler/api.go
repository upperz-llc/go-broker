package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mochi-co/mqtt/v2"
	"github.com/upperz-llc/go-broker/internal/domain"
)

type Handler struct {
	Server *mqtt.Server
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {

	client, found := h.Server.Clients.Get("test-client")
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := domain.OnlineStatusGETResponse{
		Connected: client.Closed(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
	return

}
