package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	mch "github.com/dgduncan/mochi-cloud-hooks"
	mqtt "github.com/mochi-mqtt/server/v2"

	"github.com/mochi-mqtt/server/v2/hooks/debug"
	"github.com/mochi-mqtt/server/v2/hooks/storage/redis"
	"github.com/mochi-mqtt/server/v2/listeners"
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

	// Create the new MQTT Server.
	server := mqtt.New(&mqtt.Options{})

	// CONFIGS
	ah := new(mch.HTTPAuthHook)
	gcsmh := new(mch.SecretManagerAuthHook)
	gcph := new(mch.PubsubMessagingHook)
	// *************************************

	redisConfig, err := hooks.NewRedisPersistanceHookConfig(ctx)
	if err != nil {
		server.Log.Error("", err)
		return
	}

	gcphConfig, err := hooks.NewMochiCloudHooksSecretManagerConfig(ctx)
	if err != nil {
		server.Log.Error("", err)
		return
	}

	httpauthconfig, err := hooks.NewMochiCloudHooksHTTPAuthConfig(ctx)
	if err != nil {
		server.Log.Error("", err)
		return
	}

	pshConfig, err := hooks.NewMochiCloudHooksPubSubConfig(ctx)
	if err != nil {
		server.Log.Error("", err)
		return
	}

	if err := server.AddHook(new(debug.Hook), &debug.Options{
		// ShowPacketData: true,
	}); err != nil {
		server.Log.Error("", err)
		return
	}
	if err = server.AddHook(new(redis.Hook), redisConfig); err != nil {
		log.Fatal(err)
	}
	if err = server.AddHook(gcsmh, *gcphConfig); err != nil {
		server.Log.Error("", err)
		return
	}
	if err = server.AddHook(ah, *httpauthconfig); err != nil {
		server.Log.Error("", err)
		return
	}
	if err = server.AddHook(gcph, *pshConfig); err != nil {
		server.Log.Error("", err)
		return
	}

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", nil)

	// Create a healthcheck listener
	hc := listeners.NewHTTPHealthCheck("healthcheck", ":8080", nil)

	err = server.AddListener(tcp)
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
