package msg

import (
	"context"
	"log/slog"

	"cloud.google.com/go/pubsub"
	vkit "cloud.google.com/go/pubsub/apiv1"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"

	"cspage/pkg/config"
)

func ClosePubsubClient(client *pubsub.Client) {
	slog.Debug("Closing pubsub client...")
	if client != nil {
		if err := client.Close(); err != nil {
			slog.Error("Could not close pubsub client", "err", err)
		}
	}
}

func NewPubsubClient(ctx context.Context, projectID string, clientConfig *pubsub.ClientConfig) *pubsub.Client {
	if projectID == "" {
		slog.Warn("Empty Google Project ID => disabling pubsub messages")
		return nil
	}

	var opts []option.ClientOption
	if gac := getGoogleApplicationCredentials(); gac != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(gac)))
	}

	client, err := pubsub.NewClientWithConfig(ctx, projectID, clientConfig, opts...)
	if err != nil {
		config.Die("Could not initialize pubsub client", "err", err)
	}

	return client
}

func NewPubsubConfig(pubCfg *PubsubPublisherConfig, subCfg *PubsubSubscriberConfig) *pubsub.ClientConfig {
	cfg := &pubsub.ClientConfig{}

	if pubCfg != nil {
		cfg.PublisherCallOptions = &vkit.PublisherCallOptions{
			Publish: []gax.CallOption{
				gax.WithRetry(func() gax.Retryer {
					return gax.OnCodes([]codes.Code{
						codes.Aborted,
						codes.Canceled,
						codes.Internal,
						codes.ResourceExhausted,
						codes.Unknown,
						codes.Unavailable,
						codes.DeadlineExceeded,
					}, gax.Backoff{
						Initial:    pubCfg.PubsubRetryIniBackoff,
						Max:        pubCfg.PubsubRetryMaxBackoff,
						Multiplier: pubCfg.PubsubRetryMultiplier,
					})
				}),
			},
		}
	}

	if subCfg != nil {
		_ = 1 // Not implemented because all settings are configured by the server (infra)
	}

	return cfg
}
