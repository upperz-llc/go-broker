package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/upperz-llc/go-broker/internal/hooks"
	"github.com/upperz-llc/go-broker/internal/ps"
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

	logger := client.Logger(logName).StandardLogger(logging.Info)

	// ****************** CONFIGURE SSL ****************
	certFile, err := os.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/cert.pem")
	if err != nil {
		logger.Println(err)
		return
	}

	privateKey, err := os.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/privkey.pem")
	if err != nil {
		logger.Println(err)
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
	psclient, err := pubsub.NewClient(ctx, "freezer-monitor-dev-e7d4c")
	if err != nil {
		panic(fmt.Errorf("pubsub.NewClient: %v", err))
	}
	defer client.Close()

	topic := psclient.Topic("test")
	topic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}
	topic.PublishSettings = pubsub.PublishSettings{}

	bps := ps.BrokerPubSub{
		Logger: *logger,
		Topic:  topic,
	}

	// *************************************

	// Create the new MQTT Server.
	server := mqtt.New(nil)

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)
	// _ = server.AddHook(new(hooks.FirestoreAuthHook), nil)

	examplehook := new(hooks.ExampleHook)
	examplehook.Pubsub = bps
	examplehook.Logger = logger
	_ = server.AddHook(examplehook, nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", &listeners.Config{
		TLSConfig: tlsConfig,
	})
	// tcp := listeners.NewTCP("t1", ":1883", nil)
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}

	<-done
	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")

}
