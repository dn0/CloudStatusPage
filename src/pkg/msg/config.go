//nolint:lll,tagalign // Tags of config params are manually aligned.
package msg

import (
	"os"
	"time"
)

const (
	//nolint:gosec // No secret material here.
	googleApplicationCredentialsJSON = "GOOGLE_APPLICATION_CREDENTIALS_JSON"
)

type PubsubPublisherConfig struct {
	PubsubProjectID       string        `param:"pubsub-project-id"        default:""      desc:"Google Cloud Project ID for Pub/Sub messages"`
	PubsubPublishTimeout  time.Duration `param:"pubsub-publish-timeout"   default:"10s"   desc:"timeout for publishing a Pub/Sub message"`
	PubsubRetryIniBackoff time.Duration `param:"pubsub-retry-ini-backoff" default:"250ms" desc:"initial retry period for publishing Pub/Sub messages"`
	PubsubRetryMaxBackoff time.Duration `param:"pubsub-retry-max-backoff" default:"20s"   desc:"maximum retry period for publishing Pub/Sub messages"`
	PubsubRetryMultiplier float64       `param:"pubsub-retry-multiplier"  default:"2.0"   desc:"factor by which the retry period for publishing Pub/Sub messages increases"`
}

type PubsubSubscriberConfig struct {
	PubsubProjectID string `param:"pubsub-project-id"      default:""    desc:"Google Cloud Project ID for Pub/Sub messages"`
}

func getGoogleApplicationCredentials() string {
	if value, ok := os.LookupEnv(googleApplicationCredentialsJSON); ok {
		return value
	}
	return ""
}
