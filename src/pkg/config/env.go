package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func getEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key, fallback string) bool {
	valEnv := getEnvStr(key, fallback)
	val, err := strconv.ParseBool(valEnv)
	if err != nil {
		Die(invalidConfigValue, key, valEnv, "err", err)
	}
	return val
}

func getEnvInt(key, fallback string) int64 {
	valEnv := getEnvStr(key, fallback)
	val, err := strconv.ParseInt(valEnv, 10, 64)
	if err != nil {
		Die(invalidConfigValue, key, valEnv, "err", err)
	}
	return val
}

func getEnvFloat(key, fallback string) float64 {
	valEnv := getEnvStr(key, fallback)
	val, err := strconv.ParseFloat(valEnv, 64)
	if err != nil {
		Die(invalidConfigValue, key, valEnv, "err", err)
	}
	return val
}

func getEnvDur(key, fallback string) time.Duration {
	valEnv := getEnvStr(key, fallback)
	val, err := time.ParseDuration(valEnv)
	if err != nil {
		Die(invalidConfigValue, key, valEnv, "err", err)
	}
	return val
}

func getEnvSliceStr(key, fallback string) []string {
	valEnv := getEnvStr(key, fallback)
	if valEnv == "" {
		return []string{}
	}
	return strings.Split(valEnv, stringSliceSeparator)
}
