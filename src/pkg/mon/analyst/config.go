//nolint:lll,tagalign // Tags of config params are manually aligned.
package analyst

import (
	"time"

	"cspage/pkg/config"
	"cspage/pkg/db"
	"cspage/pkg/http"
	"cspage/pkg/msg"
)

type env struct {
	config.BaseEnv

	Clouds []string `json:"clouds"`
}

// Config defaults must be set according to agent.Config's defaults.
type Config struct {
	config.BaseConfig         `param:"config-*"`
	http.HTTPConfig           `param:"http-*"`
	db.DatabaseConfig         `param:"database-*"`
	msg.PubsubPublisherConfig `param:"pubsub-*"`

	Env              env      `param:"env-*"`
	Clouds           []string `param:"cloud"              default:""           desc:"clouds for which to perform periodic analysis (can be specified multiple times)"`
	PubsubAlertTopic string   `param:"pubsub-alert-topic" default:"mon-alerts" desc:"Pub/Sub topic for ALERT messages"`

	PingCheckInterval time.Duration `param:"ping-check-interval" default:"15s"  desc:"how often to run analysis of missing pings"`
	PingLagWindow     time.Duration `param:"ping-lag-window"     default:"15m"  desc:"time window for which to look for missing pings"`
	PingLagThreshold  time.Duration `param:"ping-lag-threshold"  default:"60s"  desc:"how long an agent is considered alive after each ping"`

	ProbeCheckInterval               time.Duration `param:"probe-check-interval"       default:"30s"  desc:"how often to run analysis of probe latencies"`
	ProbeCheckWindow                 time.Duration `param:"probe-check-window"         default:"30m"  desc:"time window for which to analyze failures and latencies of standard probes"`
	ProbeLongCheckWindow             time.Duration `param:"probe-long-check-window"    default:"120m" desc:"time window for which to analyze failures and latencies of long interval probes"`
	ProbeFailureThreshold            int64         `param:"probe-failure-threshold"    default:"1"    desc:"minimum number of probe failures to trigger an alert"`
	ProbeZscoreWindow                time.Duration `param:"probe-zscore-window"        default:"8h"   desc:"time window for which to calculate z-scores of standard probes"`
	ProbeLongZscoreWindow            time.Duration `param:"probe-long-zscore-window"   default:"24h"  desc:"time window for which to calculate z-scores of long interval probes"`
	ProbeZscoreThresholdOn           float64       `param:"probe-zscore-threshold-on"  default:"3.0"  desc:"standard score (number of standard deviations) above which a probe latency is considered an issue"`
	ProbeZscoreThresholdOff          float64       `param:"probe-zscore-threshold-off" default:"2.0"  desc:"standard score (number of standard deviations) below which a probe latency is no longer considered an issue"`
	ProbeDistanceFromAvgDivThreshold int64         `param:"probe-distance-from-avg-div-threshold" default:"3"   desc:"the divisor in 'probe.DistanceFromAvg > probe.LatencyAvg / divisor' which when true the relevant probe latency can be considered an issue"`
	ProbeDistanceFromAvgMinThreshold time.Duration `param:"probe-distance-from-avg-min-threshold" default:"3ms" desc:"the minimum in 'probe.DistanceFromAvg > minimum' which when true the relevant probe latency can be considered an issue"`
	ProbeAlertSilenceAfterAgentStart time.Duration `param:"probe-alert-silence-after-agent-start" default:"2h"  desc:"time period after agent's start during which no probe alert will be created"`

	IncidentCheckInterval time.Duration `param:"incident-check-interval" default:"120s" desc:"how often to run analysis of recent alerts"`
	IncidentCheckWindow   time.Duration `param:"incident-check-window"   default:"120m" desc:"time window for which to analyze recent alerts"`
}

//nolint:gochecknoglobals // Common function.
//goland:noinspection GoUnnecessarilyExportedIdentifiers
var Die = config.Die

//nolint:gochecknoglobals // Common function.
//goland:noinspection GoUnnecessarilyExportedIdentifiers
var DieLog = config.DieLog

func NewConfig() *Config {
	cfg := &Config{
		Env: env{
			BaseEnv: config.NewBaseEnv(),
		},
	}
	config.InitConfig(cfg, &cfg.BaseConfig)
	cfg.Env.Clouds = cfg.Clouds

	return cfg
}
