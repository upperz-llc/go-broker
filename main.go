package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/logging"

	mch "github.com/dgduncan/mochi-cloud-hooks"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/upperz-llc/go-broker/internal/hooks"
	"github.com/upperz-llc/go-broker/internal/webserver"
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

	// ****************** CONFIGURE LOGGING ************

	// Creates a client.
	client, err := logging.NewClient(ctx, "freezer-monitor-dev-e7d4c")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the log to write to.
	logName := "my-log"

	logger := client.Logger(logName)

	// ****************** CONFIGURE SSL ****************
	certFile, err := os.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/cert.pem")
	if err != nil {
		logger.StandardLogger(logging.Error).Println(err)
		return
	}

	privateKey, err := os.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/privkey.pem")
	if err != nil {
		logger.StandardLogger(logging.Error).Println(err)
		return
	}

	// TLS/SSL
	cert, err := tls.X509KeyPair(certFile, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Basic TLS Config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// *************************************************

	// CONFIGS
	// TODO : Look into if this is the right way
	// fsh := new(hooks.FirestoreAuthHook)
	// fsh.Logger = logger

	ah := new(mch.HTTPAuthHook)

	gcsmh := new(mch.SecretManagerAuthHook)

	gcph := new(hooks.GCPPubsubHook)
	gcph.Logger = logger

	// *************************************

	// Create the new MQTT Server.
	server := mqtt.New(nil)

	// Allow all connections.
	// _ = server.AddHook(new(auth.AllowHook), nil)
	// _ = server.AddHook(fsh, nil)
	// _ = server.AddHook(hah, nil)

	gcphConfig, err := hooks.NewMochiCloudHooksSecretManagerConfig(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	httpauthconfig, err := hooks.NewMochiCloudHooksHTTPAuthConfig(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = server.AddHook(gcsmh, gcphConfig)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = server.AddHook(ah, httpauthconfig)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = server.AddHook(gcph, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", &listeners.Config{
		TLSConfig: tlsConfig,
	})

	// Create HTTP Stats Listener
	// stats := listeners.NewHTTPStats("stats", ":8080", nil, server.Info)
	// err = server.AddListener(stats)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go webserver.StartWebServer(server)

	<-done
	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")

}
