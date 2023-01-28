package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mochi-co/mqtt/v2"
	"github.com/upperz-llc/go-broker/internal/handler"
)

func StartWebServer(server *mqtt.Server) {
	handler := handler.Handler{
		Server: server,
	}

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/client", func(r chi.Router) {
			r.Route("/{client_id}", func(r chi.Router) {
				r.Get("/online", handler.Handle)
			})
		})
	})

	go http.ListenAndServe(":8081", r)
}
