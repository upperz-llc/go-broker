package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mch "github.com/dgduncan/mochi-cloud-hooks"
	mqtt "github.com/mochi-mqtt/server/v2"

	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/debug"
	"github.com/mochi-mqtt/server/v2/hooks/storage/badger"

	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/upperz-llc/go-broker/internal/hooks"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	ctx := context.Background()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Create the new MQTT Server.
	conf := mqtt.Options{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
	server := mqtt.New(&conf)
	server.Options.Capabilities.MaximumClientWritesPending = 32 * 1024

	// ****************** CONFIGURE LOGGING ************

	// ****************** CONFIGURE SSL ****************
	// certFile, err := os.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/cert.pem")
	// if err != nil {
	// 	server.Log.Err(err).Msg("")
	// 	return
	// }

	// privateKey, err := os.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/privkey.pem")
	// if err != nil {
	// 	server.Log.Err(err).Msg("")
	// 	return
	// }

	// // TLS/SSL
	// cert, err := tls.X509KeyPair(certFile, privateKey)
	// if err != nil {
	// 	server.Log.Err(err).Msg("")
	// 	return
	// }

	var tlsConfig *tls.Config
	if os.Getenv("FLAGS_SSL_ENABLED") == "true" {
		server.Log.Info("Configuring SSL")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
			fmt.Fprintf(w, "Hello, HTTPS world!")
		})

		leurl := "https://acme-staging-v02.api.letsencrypt.org/directory"
		if os.Getenv("LISTENERS_LETSENCRYPT_PRODUCTION") == "true" {
			leurl = autocert.DefaultACMEDirectory
		}

		// create the autocert.Manager with domains and path to the cache
		certManager := autocert.Manager{
			Client: &acme.Client{
				DirectoryURL: leurl,
			},
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(os.Getenv("LISTENERS_LETSENCRYPT_HOST")),
		}

		autocertserver := &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		log.Printf("Serving http/https for domains: %s", os.Getenv("LISTENERS_LETSENCRYPT_HOST"))
		go func() {
			// serve HTTP, which will redirect automatically to HTTPS
			h := certManager.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":http", h))
		}()

		// serve HTTPS!
		go func() {
			log.Fatal(autocertserver.ListenAndServeTLS("", ""))
		}()

		// Basic TLS Config
		tlsConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
		}
	}

	if os.Getenv("FLAGS_PUBSUB_ENABLED") == "true" {
		gcph := new(mch.PubsubMessagingHook)

		pshConfig, err := hooks.NewMochiCloudHooksPubSubConfig(ctx)
		if err != nil {
			server.Log.Error("", err)
			return
		}

		if err = server.AddHook(gcph, *pshConfig); err != nil {
			server.Log.Error("", err)
			return
		}
	}

	if os.Getenv("FLAGS_HTTP_AUTH_ENABLED") == "true" {
		ah := new(mch.HTTPAuthHook)

		httpauthconfig, err := hooks.NewMochiCloudHooksHTTPAuthConfig(ctx)
		if err != nil {
			server.Log.Error("", err)
			return
		}

		if err = server.AddHook(ah, *httpauthconfig); err != nil {
			server.Log.Error("", err)
			return
		}
	} else {
		server.AddHook(new(auth.AllowHook), nil)
	}

	if os.Getenv("FLAGS_SECRET_MANAGER_ENABLED") == "true" {
		gcsmh := new(mch.SecretManagerAuthHook)

		gcphConfig, err := hooks.NewMochiCloudHooksSecretManagerConfig(ctx)
		if err != nil {
			server.Log.Error("", err)
			return
		}

		if err = server.AddHook(gcsmh, *gcphConfig); err != nil {
			server.Log.Error("", err)
			return
		}
	}

	// CONFIGS
	// *************************************

	if os.Getenv("FLAGS_DEBUG") == "true" {
		if err := server.AddHook(new(debug.Hook), &debug.Options{
			// ShowPacketData: true,
		}); err != nil {
			server.Log.Error("", err)
			return
		}
	}

	if os.Getenv("FLAGS_STORAGE") == "true" {
		err := server.AddHook(new(badger.Hook), &badger.Options{
			Path: "/badger/.badger",
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create a TCP listener on a standard port.

	tcp := listeners.NewTCP("t1", ":1883", &listeners.Config{
		TLSConfig: tlsConfig,
	})

	// Create a healthcheck listener
	hc := listeners.NewHTTPHealthCheck("healthcheck", ":8080", nil)

	hs := listeners.NewHTTPStats("stats", ":8081", nil, server.Info)

	err := server.AddListener(tcp)
	if err != nil {
		server.Log.Error("", err)
		return
	}

	err = server.AddListener(hc)
	if err != nil {
		server.Log.Error("", err)
		return
	}

	err = server.AddListener(hs)
	if err != nil {
		server.Log.Error("", err)
		return
	}

	go func() {
		err := server.Serve()
		if err != nil {
			server.Log.Error("", err)
			return
		}
	}()

	<-done
	server.Log.Warn("caught signal, stopping...")
	server.Close()
	server.Log.Info("main.go finished")

}
