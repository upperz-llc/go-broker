package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	mqtt "github.com/mochi-mqtt/server/v2"

	"github.com/upperz-llc/go-broker/internal/handler"
)

func StartWebServer(server *mqtt.Server) {
	handler := handler.APIHandler{
		Server: server,
	}

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/client", func(r chi.Router) {
			r.Route("/{client_id}", func(r chi.Router) {
				r.Get("/connection_status", handler.HandleGetConnectionStatus)
			})
		})
	})

	go http.ListenAndServe(":80", r)
}
