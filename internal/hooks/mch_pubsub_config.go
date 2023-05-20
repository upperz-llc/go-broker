package hooks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	mch "github.com/dgduncan/mochi-cloud-hooks"
	"github.com/upperz-llc/go-broker/internal/admin"
)

func NewMochiCloudHooksPubSubConfig(ctx context.Context) (*mch.PubsubMessagingHookConfig, error) {
	// Pull Env Variables
	pid, found := os.LookupEnv("GCP_PROJECT_ID")
	if !found {
		return nil, errors.New("GCP_PROJECT_ID not found")
	}
	bt, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_PUBLISH")
	if !found {
		return nil, errors.New("BROKER_HOOK_GCPPUBSUB_TOPIC_PUBLISH not found")
	}
	st, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_SUBSCRIBE")
	if !found {
		return nil, errors.New("BROKER_HOOK_GCPPUBSUB_TOPIC_SUBSCRIBE not found")
	}
	ct, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_CONNECT")
	if !found {
		return nil, errors.New("BROKER_HOOK_GCPPUBSUB_TOPIC_CONNECT not found")
	}
	lwtt, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_LWT")
	if !found {
		return nil, errors.New("BROKER_HOOK_GCPPUBSUB_TOPIC_LWT not found")
	}
	oset, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_ONSESSIONESTABLISHED")
	if !found {
		return nil, errors.New("BROKER_HOOK_GCPPUBSUB_TOPIC_ONSESSIONESTABLISHED not found")
	}

	adminclient, err := admin.NewAdmin(ctx)
	if err != nil {
		return nil, err
	}

	disallowList := make([]string, 0)
	disallowList = append(disallowList, adminclient.GetAdminCredentials())

	// Create and configure pubsub client
	pc, err := pubsub.NewClient(ctx, pid)
	if err != nil {
		panic(fmt.Errorf("pubsub.NewClient: %v", err))
	}

	pubslishtopic := pc.Topic(bt)
	pubslishtopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	subscribetopic := pc.Topic(st)
	subscribetopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	connecttopic := pc.Topic(ct)
	connecttopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	willtopic := pc.Topic(lwtt)
	willtopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	onSessionEstablishedTopic := pc.Topic(oset)
	onSessionEstablishedTopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	return &mch.PubsubMessagingHookConfig{
		ConnectTopic:              connecttopic,
		OnSessionEstablishedTopic: onSessionEstablishedTopic,
		PublishTopic:              pubslishtopic,
		SubscribeTopic:            subscribetopic,
		WillTopic:                 willtopic,
		DisallowList:              disallowList,
	}, nil
}
