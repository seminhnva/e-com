package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallBack string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallBack
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsDuration, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return valAsDuration
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return boolVal
}
