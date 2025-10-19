package helpers

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetenvDuration(key string, defaultSeconds int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(defaultSeconds) * time.Second
	}
	seconds, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid value for %s, using default %d", key, defaultSeconds)
		return time.Duration(defaultSeconds) * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func GetenvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	return value
}
