//nolint:lll,tagalign // Tags of config params are manually aligned.
package scribe

import (
	"cspage/pkg/config"
	"cspage/pkg/db"
	"cspage/pkg/http"
	"cspage/pkg/msg"
)

type env struct {
	config.BaseEnv

	PubsubSubscriptions []string `json:"pubsub_subscriptions"`
}

type Config struct {
	config.BaseConfig          `param:"config-*"`
	http.HTTPConfig            `param:"http-*"`
	db.DatabaseConfig          `param:"database-*"`
	msg.PubsubSubscriberConfig `param:"pubsub-*"`

	Env                      env      `param:"env-*"`
	PubsubPingSubscriptions  []string `param:"pubsub-ping-subscription"  default:"" desc:"pub/sub subscription for PING and AGENT messages (can be specified multiple times)"`
	PubsubProbeSubscriptions []string `param:"pubsub-probe-subscription" default:"" desc:"Pub/Sub subscription for PROBE messages (can be specified multiple times)"`
}

func NewConfig() *Config {
	cfg := &Config{
		Env: env{
			BaseEnv: config.NewBaseEnv(),
		},
	}
	config.InitConfig(cfg, &cfg.BaseConfig)
	//nolint:gocritic // We want to store the concatenated slices in the env var.
	cfg.Env.PubsubSubscriptions = append(cfg.PubsubPingSubscriptions, cfg.PubsubProbeSubscriptions...)

	return cfg
}
