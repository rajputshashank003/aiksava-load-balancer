package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getInt(key string, fallback int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil || val <= 0 {
		return fallback
	}
	return val
}

func GetBackendCount() int {
	backendCount := getInt("BACKEND_COUNT", 1)
	fmt.Println("backend count:", backendCount)
	return backendCount
}

func GetMaxUsersPerServer() int {
	fmt.Println("Fetching MAX_USERS_PER_SERVER from environment...", os.Getenv("MAX_USERS_PER_SERVER"))
	maxUsers := getInt("MAX_USERS_PER_SERVER", 10)
	fmt.Println("Max users per server set to:", maxUsers)
	return maxUsers
}

func GetSessionTTL() time.Duration {
	ttlSeconds := getInt("SESSION_TTL_SECONDS", 600)
	fmt.Println("max seconds:", ttlSeconds)
	return time.Duration(ttlSeconds) * time.Second
}

func GetAllowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		// Default origins matching your frontend config
		return []string{
			"https://aiksava.onrender.com",
			"https://aiksava.vercel.app",
			"http://localhost:8000",
		}
	}
	// Split comma-separated origins
	var result []string
	for _, origin := range splitAndTrim(origins, ",") {
		if origin != "" {
			result = append(result, origin)
		}
	}
	return result
}

func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}
