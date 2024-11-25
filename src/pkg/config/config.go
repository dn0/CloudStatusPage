//nolint:lll,tagalign // Tags of config params are manually aligned.
package config

import (
	"time"
)

type BaseEnv struct {
	Tag      string `json:"tag"`
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
}

type BaseConfig struct {
	Debug     bool   `param:"debug"      default:"false" desc:"enable /debug/pprof/ endpoint"`
	DryRun    bool   `param:"dry-run"    default:"false" desc:"enable 'dry run' mode (behaves differently for each service)"`
	LogFormat string `param:"log-format" default:"text"  desc:"log format: text or json"`
	LogLevel  string `param:"log-level"  default:"info"  desc:"log level: debug, info, warn, error"`

	WorkerStartTimeout    time.Duration `param:"worker-start-timeout"    default:"20s"   desc:"maximum time to wait for initialization functions to complete"`
	WorkerShutdownTimeout time.Duration `param:"worker-shutdown-timeout" default:"20s"   desc:"maximum time to wait for graceful shutdown to complete"`
	WorkerStopTimeout     time.Duration `param:"worker-stop-timeout"     default:"10s"   desc:"maximum time to wait for cleanup functions to complete"`
	WorkerTaskTimeout     time.Duration `param:"worker-task-timeout"     default:"600s"  desc:"maximum allowed time for main job task function"`
}
