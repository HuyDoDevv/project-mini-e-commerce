package utils

import "os"

func GetEnv(key, defaultValue string) string {
	valueKey := os.Getenv(key)
	if valueKey == "" {
		return defaultValue
	}
	return valueKey
}
