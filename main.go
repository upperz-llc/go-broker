package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/upperz-llc/go-broker/internal/hooks"
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

	// certFile, err := ioutil.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/cert.pem")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// privateKey, err := ioutil.ReadFile("etc/letsencrypt/live/testbroker.dev.upperz.org/privkey.pem")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // TLS/SSL
	// cert, err := tls.X509KeyPair(certFile, privateKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Basic TLS Config
	// tlsConfig := &tls.Config{
	// 	Certificates: []tls.Certificate{cert},
	// }

	// // CONFIGS
	// client, err := pubsub.NewClient(ctx, "freezer-monitor-dev-e7d4c")
	// if err != nil {
	// 	panic(fmt.Errorf("pubsub.NewClient: %v", err))
	// }
	// defer client.Close()

	// topic := client.Topic("test")
	// topic.PublishSettings = pubsub.PublishSettings{
	// 	DelayThreshold: 1 * time.Second,
	// 	CountThreshold: 10,
	// }
	// topic.PublishSettings = pubsub.PublishSettings{}

	// bps := ps.BrokerPubSub{
	// 	Topic: topic,
	// }

	// *************************************

	// Create the new MQTT Server.
	server := mqtt.New(nil)

	// Allow all connections.
	// _ = server.AddHook(new(auth.AllowHook), nil)
	_ = server.AddHook(new(hooks.FirestoreAuthHook), nil)

	// examplehook := new(hooks.ExampleHook)
	// examplehook.Pubsub = bps
	// _ = server.AddHook(examplehook, nil)

	// Create a TCP listener on a standard port.
	// tcp := listeners.NewTCP("t1", ":1883", &listeners.Config{
	// 	TLSConfig: tlsConfig,
	// })
	tcp := listeners.NewTCP("t1", ":1883", nil)
	err := server.AddListener(tcp)
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
