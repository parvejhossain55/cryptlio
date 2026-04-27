package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv     string
	ServerPort string
	JWTSecret  string
	JWTExpiry  time.Duration

	// Service URLs
	AuthServiceURL     string
	UserServiceURL     string
	WalletServiceURL   string
	TradeEngineURL     string
	KYCEngineURL       string
	NotificationURL    string
	DisputeServiceURL  string
	MerchantServiceURL string

	// Rate limiting
	RateLimitEnabled  bool
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// CORS
	CorsAllowedOrigins []string
}

func Load() (*Config, error) {
	rateLimitEnabled, _ := strconv.ParseBool(getEnv("RATE_LIMIT_ENABLED", "true"))
	rateLimitRequests, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1m"))

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiry:  jwtExpiry,

		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
		UserServiceURL:     getEnv("USER_SERVICE_URL", "http://localhost:8080"),
		WalletServiceURL:   getEnv("WALLET_SERVICE_URL", "http://localhost:8081"),
		TradeEngineURL:     getEnv("TRADE_ENGINE_URL", "http://localhost:8082"),
		KYCEngineURL:       getEnv("KYC_SERVICE_URL", "http://localhost:8083"),
		NotificationURL:    getEnv("NOTIFICATION_URL", "http://localhost:8084"),
		DisputeServiceURL:  getEnv("DISPUTE_SERVICE_URL", "http://localhost:8085"),
		MerchantServiceURL: getEnv("MERCHANT_SERVICE_URL", "http://localhost:8082"),

		RateLimitEnabled:   rateLimitEnabled,
		RateLimitRequests:  rateLimitRequests,
		RateLimitWindow:    rateLimitWindow,
		CorsAllowedOrigins: []string{getEnv("CORS_ALLOWED_ORIGINS", "*")},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
