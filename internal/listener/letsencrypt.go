package listener

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/rs/zerolog"
)

// LetsEncrypt is a listener for providing an HTTP healthcheck endpoint.
type LetsEncrypt struct {
	sync.RWMutex
	id      string            // the internal id of the listener
	address string            // the network address to bind to
	config  *listeners.Config // configuration values for the listener
	listen  *http.Server      // the http server
	log     *zerolog.Logger   // server logger
	end     uint32            // ensure the close methods are only called once
}

// NewLetsEncrypt initialises and returns a new HTTP listener, listening on an address.
func NewLetsEncrypt(id, address string, config *listeners.Config) *LetsEncrypt {
	if config == nil {
		config = new(listeners.Config)
	}
	return &LetsEncrypt{
		id:      id,
		address: address,
		config:  config,
	}
}

// ID returns the id of the listener.
func (l *LetsEncrypt) ID() string {
	return l.id
}

// Address returns the address of the listener.
func (l *LetsEncrypt) Address() string {
	return l.address
}

// Protocol returns the address of the listener.
func (l *LetsEncrypt) Protocol() string {
	if l.listen != nil && l.listen.TLSConfig != nil {
		return "https"
	}

	return "http"
}

// Init initializes the listener.
func (l *LetsEncrypt) Init(log *zerolog.Logger) error {
	l.log = log

	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
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
func (l *LetsEncrypt) Serve(establish listeners.EstablishFn) {
	if l.listen.TLSConfig != nil {
		l.listen.ListenAndServeTLS("", "")
	} else {
		l.listen.ListenAndServe()
	}
}

// Close closes the listener and any client connections.
func (l *LetsEncrypt) Close(closeClients listeners.CloseFn) {
	l.Lock()
	defer l.Unlock()

	if atomic.CompareAndSwapUint32(&l.end, 0, 1) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		l.listen.Shutdown(ctx)
	}

	closeClients(l.id)
}
