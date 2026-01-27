// Package env provides utilities to work with environment variables
package env

import (
	"strings"
	"strconv"
	"os"
	"fmt"
)

func GetString (key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetStringRequired (key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("environment variable %s is required", key))
}

func GetInt (key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(fmt.Errorf("failed to convert env variable %s to int: %w", key, err))
		}
		return intValue
	}
	return defaultValue
}

func GetIntRequired (key string) int {
	value := GetStringRequired(key)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Errorf("failed to convert env variable %s to int: %w", key, err))
	}
	return intValue
}

func GetBool (key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		value = strings.ToLower(strings.TrimSpace(value))
		switch value {
		case "true", "1", "yes", "y", "on", "enable":
			return true
		case "false", "0", "no", "n", "off", "disable":
			return false
		default:
			panic(fmt.Sprintf("boolean env variable %s must be true/false, 1/0, yes/no or on/off", key))
		}
	}
	return defaultValue
}

func GetBoolRequired (key string) bool {
	value := GetStringRequired(key)
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "true", "1", "yes", "y", "on", "enable":
		return true
	case "false", "0", "no", "n", "off", "disable":
	    return false
	default:
		panic(fmt.Sprintf("boolean env variable %s must be true/false, 1/0, yes/no or on/off", key))
	}
}

