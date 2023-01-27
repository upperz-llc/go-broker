package handler

import (
	"net/http"

	"github.com/mochi-co/mqtt/v2"
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

	if client.Closed() {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}
