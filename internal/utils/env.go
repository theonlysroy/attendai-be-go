package utils

import "os"

func Getkey(key string, defaultVal string) string {
	if key == "" {
		return defaultVal
	}
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
