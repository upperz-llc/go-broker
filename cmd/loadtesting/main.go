package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/hooks/debug"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"
	"github.com/upperz-llc/go-broker/internal/webserver"
)

func main() {
	// ctx := context.Background()

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
	// pid, found := os.LookupEnv("GCP_PROJECT_ID")
	// if !found {
	// 	log.Fatal("GCP_PROJECT_ID not found")
	// }

	// Create GCP Zap Logger
	// gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, pid, "mochi-broker", zlg.CloudLoggingOptions{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// gcpZeroLogger := zerolog.New(gcpWriter)
	// debugLogger := gcpZeroLogger.Level(zerolog.DebugLevel)
	// server.Log = &debugLogger

	_ = server.AddHook(new(debug.Hook), &debug.Options{})
	_ = server.AddHook(new(auth.AllowHook), nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", nil)

	// Create a healthcheck listener
	hc := listeners.NewHTTPHealthCheck("healthcheck", ":8080", nil)

	err := server.AddListener(tcp)
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
