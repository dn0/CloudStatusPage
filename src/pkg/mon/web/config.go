//nolint:lll,tagalign // Tags of config params are manually aligned.
package web

import (
	"cspage/pkg/config"
	"cspage/pkg/db"
	"cspage/pkg/http"
)

type env struct {
	config.BaseEnv
}

type Config struct {
	config.BaseConfig `param:"config-*"`
	http.HTTPConfig   `param:"http-*"`
	db.DatabaseConfig `param:"database-*"`

	Env                   env    `param:"env-*"`
	HTTPBasicAuthRealm    string `param:"http-basic-auth-realm"    default:"SecretArea" desc:"Simple protection of secret endpoints"`
	HTTPBasicAuthUsername string `param:"http-basic-auth-username" default:"cloudstatus" desc:"Simple protection of secret endpoints"`
	HTTPBasicAuthPassword string `param:"http-basic-auth-password" default:"nbusr123" desc:"Simple protection of secret endpoints"`
}

func NewConfig() *Config {
	cfg := &Config{
		Env: env{
			BaseEnv: config.NewBaseEnv(),
		},
	}
	config.InitConfig(cfg, &cfg.BaseConfig)

	return cfg
}
