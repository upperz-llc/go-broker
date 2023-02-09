package hooks

import (
	"bytes"
	"context"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	httpauth "github.com/upperz-llc/http-auth-backend/pkg/api"
)

type HTTPAuthHook struct {
	HTTPClient *httpauth.HTTPAuthBackendClient
	Logger     *logging.Logger
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
	httpclient, err := httpauth.NewClient(context.Background())
	if err != nil {
		return err
	}

	h.HTTPClient = httpclient
	h.Logger.StandardLogger(logging.Debug).Println("initialized httpauth")
	return nil
}

func (h *HTTPAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	allowed, err := h.HTTPClient.CheckClientAuth(context.Background(), cl.ID, string(pk.Connect.Username), string(pk.Connect.Password))
	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}
	return allowed
}

func (h *HTTPAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	allowed, err := h.HTTPClient.CheckClientACLs(context.Background(), cl.ID, string(cl.Properties.Username), topic, strconv.FormatBool(write))
	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}

	return allowed
}
