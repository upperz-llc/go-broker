package hooks

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
)

type HTTPAuthHook struct {
	Client *api.HTTPAuthBackendClient
	Logger *logging.Logger
	mqtt.HookBase
}

func (h *HTTPAuthHook) ID() string {
	return "http-auth-hook"
}

func (h *HTTPAuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnACLCheck,
		mqtt.OnConnectAuthenticate,
	}, []byte{b})
}

func (h *HTTPAuthHook) Init(config any) error {
	httpauthclient, err := api.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	h.Client = httpauthclient

	return nil
}

func (h *HTTPAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	allowed, err := h.Client.CheckClientAuth(context.Background(), cl.ID, string(pk.Connect.Username), string(pk.Connect.Password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return allowed
}

func (h *HTTPAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	allowed, err := h.Client.CheckClientACLs(context.Background(), cl.ID, string(cl.Properties.Username), topic, strconv.FormatBool(write))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return allowed
}
