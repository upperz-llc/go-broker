package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mch "github.com/dgduncan/mochi-cloud-hooks"
	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/debug"
	"github.com/mochi-co/mqtt/v2/hooks/storage/redis"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"
	"github.com/upperz-llc/go-broker/internal/hooks"
	"github.com/upperz-llc/go-broker/internal/webserver"
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
	server := mqtt.New(&mqtt.Options{})
	l := server.Log.Level(zerolog.DebugLevel)
	server.Log = &l

	// ****************** CONFIGURE LOGGING ************

	// pull project id from env
	pid, found := os.LookupEnv("GCP_PROJECT_ID")
	if !found {
		log.Fatal("GCP_PROJECT_ID not found")
	}

	// Creates gcp cloud logger client.
	// client, err := logging.NewClient(ctx, pid)
	// if err != nil {
	// 	log.Fatalf("Failed to create client: %v", err)
	// }
	// defer client.Close()

	// Sets the name of the log to write to.
	// logName := "mochi-broker"

	// logger := client.Logger(logName)

	// Create GCP Zap Logger
	gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, pid, "mochi-broker", zlg.CloudLoggingOptions{})
	if err != nil {
		log.Fatal(err)
	}

	gcpZeroLogger := zerolog.New(gcpWriter)
	debugLogger := gcpZeroLogger.Level(zerolog.DebugLevel)
	server.Log = &debugLogger

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
		fmt.Fprintf(w, "Hello, HTTPS world!")
	})

	// create the autocert.Manager with domains and path to the cache
	certManager := autocert.Manager{
		Client: &acme.Client{
			DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
		},
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("autocert.dev.upperz.org"),
	}

	autocertserver := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	log.Printf("Serving http/https for domains: %s", "autocert.dev.upperz.org")
	go func() {
		// serve HTTP, which will redirect automatically to HTTPS
		h := certManager.HTTPHandler(nil)
		log.Fatal(http.ListenAndServe(":http", h))
	}()

	// serve HTTPS!
	go func() {
		fmt.Println(autocertserver.Addr)
		log.Fatal(autocertserver.ListenAndServeTLS("", ""))
	}()

	// *************************************************

	// Basic TLS Config
	tlsConfig := &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}

	// CONFIGS
	ah := new(mch.HTTPAuthHook)
	gcsmh := new(mch.SecretManagerAuthHook)
	gcph := new(mch.PubsubMessagingHook)
	// *************************************

	redisConfig, err := hooks.NewRedisPersistanceHookConfig(ctx)
	if err != nil {
		server.Log.Err(err).Msg("")
		return
	}

	gcphConfig, err := hooks.NewMochiCloudHooksSecretManagerConfig(ctx)
	if err != nil {
		server.Log.Err(err).Msg("")
		return
	}

	httpauthconfig, err := hooks.NewMochiCloudHooksHTTPAuthConfig(ctx)
	if err != nil {
		server.Log.Err(err).Msg("")
		return
	}

	pshConfig, err := hooks.NewMochiCloudHooksPubSubConfig(ctx)
	if err != nil {
		server.Log.Err(err).Msg("")
		return
	}

	if err := server.AddHook(new(debug.Hook), &debug.Options{
		// ShowPacketData: true,
	}); err != nil {
		server.Log.Err(err).Msg("")
		return
	}
	if err = server.AddHook(new(redis.Hook), redisConfig); err != nil {
		log.Fatal(err)
	}
	if err = server.AddHook(gcsmh, *gcphConfig); err != nil {
		server.Log.Err(err).Msg("")
		return
	}
	if err = server.AddHook(ah, *httpauthconfig); err != nil {
		server.Log.Err(err).Msg("")
		return
	}
	if err = server.AddHook(gcph, *pshConfig); err != nil {
		server.Log.Err(err).Msg("")
		return
	}

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", &listeners.Config{
		TLSConfig: tlsConfig,
	})

	// Create a healthcheck listener
	hc := listeners.NewHTTPHealthCheck("healthcheck", ":8080", nil)

	err = server.AddListener(tcp)
	if err != nil {
		server.Log.Err(err).Msg("")
		return
	}

	err = server.AddListener(hc)
	if err != nil {
		server.Log.Err(err).Msg("")
		return
	}

	go func() {
		err := server.Serve()
		if err != nil {
			server.Log.Err(err).Msg("")
			return
		}
	}()

	go webserver.StartWebServer(server)

	<-done
	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")

}
