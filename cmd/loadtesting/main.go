package main

import (
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"

	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/debug"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/upperz-llc/go-broker/internal/webserver"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Create the new MQTT Server.
	server := mqtt.New(&mqtt.Options{})

	// ****************** CONFIGURE LOGGING ************

	_ = server.AddHook(new(debug.Hook), &debug.Options{})
	_ = server.AddHook(new(auth.AllowHook), nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", nil)

	// Create a healthcheck listener
	hc := listeners.NewHTTPHealthCheck("healthcheck", ":8080", nil)

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

	go func() {
		err := server.Serve()
		if err != nil {
			server.Log.Error("", err)
			return
		}
	}()

	go webserver.StartWebServer(server)

	<-done
	server.Log.Warn("caught signal, stopping...")
	server.Close()
	server.Log.Info("main.go finished")

}
