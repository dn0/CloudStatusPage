//nolint:lll,tagalign // Tags of config params are manually aligned.
package http

import (
	"time"
)

//nolint:revive // Name is OK as it's used for embedding by other Configs.
//goland:noinspection GoNameStartsWithPackageName
type HTTPConfig struct {
	HTTPListenAddr       string        `param:"http-listen-addr"       default:":8000" desc:"HTTP server's listen address"`
	HTTPShutdownTimeout  time.Duration `param:"http-shutdown-timeout"  default:"20s"   desc:"HTTP server's maximum time to wait for graceful shutdown to complete"`
	HTTPReadTimeout      time.Duration `param:"http-read-timeout"      default:"5s"    desc:"HTTP server's request read timeout"`
	HTTPWriteTimeout     time.Duration `param:"http-write-timeout"     default:"10s"   desc:"HTTP server's request write timeout"`
	HTTPIdleTimeout      time.Duration `param:"http-idle-timeout"      default:"15s"   desc:"HTTP server's maximum wait time for the next request when keep-alives are enabled"`
	HTTPMaxHeaderSize    int64         `param:"http-max-header-size"   default:"2048"  desc:"HTTP server's maximum bytes to read during header parsing"`
	HTTPMaxBodySize      int64         `param:"http-max-body-size"     default:"16384" desc:"HTTP server's maximum bytes to read during request body reading"`
	HTTPCompressionLevel int64         `param:"http-compression-level" default:"0"     desc:"HTTP server's compression level, 0 = disabled"`
}
