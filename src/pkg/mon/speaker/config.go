//nolint:tagalign // Tags of config params are manually aligned.
package speaker

import (
	"cspage/pkg/config"
	"cspage/pkg/db"
	"cspage/pkg/mon/analyst"
)

// Config for speaker is a special one. The embedded configs are pointers without 'param' and set from the analyst config.
type Config struct {
	*config.BaseConfig
	*db.DatabaseConfig

	SiteURL             string `param:"site-url"              default:"https://cloudstatus.page"`
	TwitterAPIKey       string `param:"twitter-api-key"       default:""`
	TwitterAPISecret    string `param:"twitter-api-secret"    default:""`
	TwitterAccessToken  string `param:"twitter-access-token"  default:""`
	TwitterAccessSecret string `param:"twitter-access-secret" default:""`
}

func NewConfigFromAnalyst(analystCfg *analyst.Config) *Config {
	cfg := &Config{
		BaseConfig:     &analystCfg.BaseConfig,
		DatabaseConfig: &analystCfg.DatabaseConfig,
	}
	config.InitConfig(cfg, cfg.BaseConfig)

	return cfg
}
