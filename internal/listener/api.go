package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/upperz-llc/go-broker/pkg/api"
)

// API is a listener for providing an HTTP healthcheck endpoint.
type API struct {
	sync.RWMutex
	id           string            // the internal id of the listener
	address      string            // the network address to bind to
	config       *listeners.Config // configuration values for the listener
	listen       *http.Server      // the http server
	log          *slog.Logger      // server logger
	end          uint32            // ensure the close methods are only called once
	inlineClient *mqtt.Server      // the mqtt server used for inline client
}

// NewLetsEncrypt initialises and returns a new HTTP listener, listening on an address.
func NewAPI(id, address string, inlineClient *mqtt.Server, config *listeners.Config) *API {
	if config == nil {
		config = new(listeners.Config)
	}
	return &API{
		id:           id,
		address:      address,
		config:       config,
		inlineClient: inlineClient,
	}
}

// ID returns the id of the listener.
func (l *API) ID() string {
	return l.id
}

// Address returns the address of the listener.
func (l *API) Address() string {
	return l.address
}

// Protocol returns the address of the listener.
func (l *API) Protocol() string {
	if l.listen != nil && l.listen.TLSConfig != nil {
		return "https"
	}

	return "http"
}

// Init initializes the listener.
func (l *API) Init(log *slog.Logger) error {
	l.log = log

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/message", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var msgReq api.MessageRequest
		if err := json.Unmarshal(body, &msgReq); err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}

		fmt.Println(msgReq.Topic, msgReq.Payload, msgReq.Retain, msgReq.QoS)

		if err := l.inlineClient.Publish(msgReq.Topic, msgReq.Payload, msgReq.Retain, msgReq.QoS); err != nil {
			http.Error(w, "Error publishing message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	})
	l.listen = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         l.address,
		Handler:      mux,
	}

	if l.config.TLSConfig != nil {
		l.listen.TLSConfig = l.config.TLSConfig
	}

	return nil
}

// Serve starts listening for new connections and serving responses.
func (l *API) Serve(establish listeners.EstablishFn) {
	if l.listen.TLSConfig != nil {
		l.listen.ListenAndServeTLS("", "")
	} else {
		l.listen.ListenAndServe()
	}
}

// Close closes the listener and any client connections.
func (l *API) Close(closeClients listeners.CloseFn) {
	l.Lock()
	defer l.Unlock()

	if atomic.CompareAndSwapUint32(&l.end, 0, 1) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		l.listen.Shutdown(ctx)
	}

	closeClients(l.id)
}
