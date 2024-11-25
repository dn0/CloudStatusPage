package config

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"time"
)

const invalidConfigValue = "Invalid config value"

//nolint:gochecknoglobals // This is a constant set during build time (see Makefile)
var Version string

func NewBaseEnv() BaseEnv {
	hostname, err := os.Hostname()
	if err != nil {
		Die("Could not fetch OS hostname", "err", err)
	}

	return BaseEnv{
		Tag:      Tag,
		Version:  Version,
		Hostname: hostname,
	}
}

func InitConfig(cfg any, baseCfg *BaseConfig) {
	initConfigFlags(reflect.ValueOf(cfg).Elem(), "")
	flag.Parse()
	initLog(baseCfg.LogLevel, baseCfg.LogFormat)
}

func InitConfigFlags(cfg any, flagPrefix string) {
	initConfigFlags(reflect.ValueOf(cfg).Elem(), flagPrefix)
}

//nolint:cyclop // Higher complexity is unavoidable here.
func initConfigFlags(reflection reflect.Value, flagPrefix string) {
	for i := 0; i < reflection.NumField(); i++ {
		fieldT := reflection.Type().Field(i)
		fParam := fieldT.Tag.Get("param")
		if fParam == "" {
			continue
		}
		fParam = flagPrefix + fParam
		fEnvParam := strings.ToUpper(strings.ReplaceAll(fParam, "-", "_"))
		fDesc := fieldT.Tag.Get("desc")
		fDefault := fieldT.Tag.Get("default")
		fieldV := reflection.Field(i)

		//nolint:exhaustive // This is OK as we control the config types.
		switch fieldV.Kind() {
		case reflect.Struct:
			initConfigFlags(fieldV, flagPrefix+fieldT.Tag.Get("prefix"))
		case reflect.Bool, reflect.Int64, reflect.Float64, reflect.String, reflect.Slice:
			switch fPointer := fieldV.Addr().Interface().(type) {
			case *bool:
				flag.BoolVar(fPointer, fParam, getEnvBool(fEnvParam, fDefault), fDesc)
			case *string:
				flag.StringVar(fPointer, fParam, getEnvStr(fEnvParam, fDefault), fDesc)
			case *int64:
				flag.Int64Var(fPointer, fParam, getEnvInt(fEnvParam, fDefault), fDesc)
			case *float64:
				flag.Float64Var(fPointer, fParam, getEnvFloat(fEnvParam, fDefault), fDesc)
			case *time.Duration:
				flag.DurationVar(fPointer, fParam, getEnvDur(fEnvParam, fDefault), fDesc)
			case *[]string:
				flag.Var(newStringSliceValue(getEnvSliceStr(fEnvParam, fDefault), fPointer), fParam, fDesc)
			}
		default:
			panic("unhandled config type")
		}
	}
}
