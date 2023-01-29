package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mochi-co/mqtt/v2"
	"github.com/upperz-llc/go-broker/pkg/domain"
)

type APIHandler struct {
	Server *mqtt.Server
}

func (ah *APIHandler) HandleGetConnectionStatus(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "client_id")

	// get client from MQTT server
	client, found := ah.Server.Clients.Get(clientID)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := domain.ConnectionStatusGETResponse{
		Connected: !client.Closed(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
