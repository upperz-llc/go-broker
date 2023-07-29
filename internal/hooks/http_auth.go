package hooks

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/upperz-llc/go-broker/internal/admin"
	"github.com/upperz-llc/go-broker/internal/httpauth"
)

type HTTPAuthHook struct {
	admin      *admin.Admin
	httpClient *httpauth.HTTPAuthBackendClient
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
	ctx := context.Background()

	admin, err := admin.NewAdmin(ctx)
	if err != nil {
		return err
	}

	httpclient, err := httpauth.NewClient(ctx, nil)
	if err != nil {
		return err
	}

	h.httpClient = httpclient
	h.admin = admin

	h.Logger.StandardLogger(logging.Debug).Println("initialized http-auth-hook")
	return nil
}

func (h *HTTPAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {

	fmt.Println(h.admin.GetAdminCredentials())
	// CHECK ADMIN
	if string(pk.Connect.Username) == h.admin.GetAdminCredentials() {
		return true
	}
	// ****************************

	// Call HTTP auth backend
	allowed, err := h.httpClient.CheckClientAuth(context.Background(), cl.ID, string(pk.Connect.Username), string(pk.Connect.Password))
	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}
	// ****************************

	return allowed
}

func (h *HTTPAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	// CHECK ADMIN
	if string(cl.Properties.Username) == h.admin.GetAdminCredentials() {
		return true
	}
	// ****************************

	// Call HTTP auth backend
	allowed, err := h.httpClient.CheckClientACLs(context.Background(), cl.ID, string(cl.Properties.Username), topic, strconv.FormatBool(write))
	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}
	// ****************************

	return allowed
}
