package router

import (
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/cryplio/services/api-gateway/internal/config"
	"github.com/cryplio/services/api-gateway/internal/handlers"
	"github.com/cryplio/services/api-gateway/internal/middleware"
	"github.com/cryplio/services/api-gateway/internal/proxy"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if cfg.AppEnv == "development" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())

	// Health check endpoints (public)
	r.GET("/health", handlers.HealthCheckHandler())
	r.GET("/live", handlers.LivenessHandler())
	r.GET("/ready", handlers.ReadinessHandler())

	// Check if auth service is external or local (development mode)
	authURL, _ := url.Parse(cfg.AuthServiceURL)
	isLocalAuth := authURL.Host == "localhost:8080" || authURL.Host == "127.0.0.1:8080" || authURL.Host == ""

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes (no auth required)
		public := v1.Group("/")
		{
			// Public ads listing
			public.GET("/ads", proxy.ProxyHandler(cfg.TradeEngineURL))
			public.GET("/market/rates", proxy.ProxyHandler(cfg.TradeEngineURL))
		}

		// Auth routes: use local handlers in dev, proxy otherwise
		if isLocalAuth {
			public.POST("/auth/register", handlers.RegisterHandler(cfg.JWTSecret, cfg.JWTExpiry))
			public.POST("/auth/login", handlers.LoginHandler(cfg.JWTSecret, cfg.JWTExpiry))
		} else {
			public.POST("/auth/register", proxy.ProxyHandler(cfg.AuthServiceURL))
			public.POST("/auth/login", proxy.ProxyHandler(cfg.AuthServiceURL))
		}

		// Authenticated routes (require JWT)
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// User routes (proxy to user service or use local mock)
			if isLocalAuth {
				auth.GET("/users/me", handlers.GetUserHandler())
				auth.PUT("/users/me", handlers.UpdateUserHandler())
			} else {
				auth.GET("/users/me", proxy.ProxyHandler(cfg.UserServiceURL))
				auth.PUT("/users/me", proxy.ProxyHandler(cfg.UserServiceURL))
			}

			// KYC
			auth.POST("/kyc/submit", proxy.ProxyHandler(cfg.KYCEngineURL))
			auth.GET("/kyc/status", proxy.ProxyHandler(cfg.KYCEngineURL))

			// Ads (authenticated)
			auth.POST("/ads", proxy.ProxyHandler(cfg.TradeEngineURL))
			auth.PUT("/ads/:id", proxy.ProxyHandler(cfg.TradeEngineURL))
			auth.DELETE("/ads/:id", proxy.ProxyHandler(cfg.TradeEngineURL))

			// Trades
			auth.POST("/trades", proxy.ProxyHandler(cfg.TradeEngineURL))
			auth.POST("/trades/:id/paid", proxy.ProxyHandler(cfg.TradeEngineURL))
			auth.POST("/trades/:id/release", proxy.ProxyHandler(cfg.TradeEngineURL))
			auth.POST("/trades/:id/cancel", proxy.ProxyHandler(cfg.TradeEngineURL))
			auth.POST("/trades/:id/dispute", proxy.ProxyHandler(cfg.DisputeServiceURL))
			auth.GET("/trades", proxy.ProxyHandler(cfg.TradeEngineURL))
			auth.GET("/trades/:id", proxy.ProxyHandler(cfg.TradeEngineURL))

			// Wallet
			auth.GET("/wallet/balance", proxy.ProxyHandler(cfg.WalletServiceURL))
			auth.POST("/wallet/withdraw", proxy.ProxyHandler(cfg.WalletServiceURL))
			auth.POST("/wallet/deposit", proxy.ProxyHandler(cfg.WalletServiceURL))
			auth.GET("/wallet/history", proxy.ProxyHandler(cfg.WalletServiceURL))

			// Notifications
			auth.GET("/notifications", proxy.ProxyHandler(cfg.NotificationURL))
			auth.PATCH("/notifications/:id/read", proxy.ProxyHandler(cfg.NotificationURL))

			// Merchant (future)
			auth.GET("/merchant/dashboard", proxy.ProxyHandler(cfg.MerchantServiceURL))
		}

		// Optional auth routes (work with or without token)
		optional := v1.Group("/")
		optional.Use(middleware.OptionalAuth(cfg.JWTSecret))
		{
			optional.GET("/ads/:id", proxy.ProxyHandler(cfg.TradeEngineURL))
		}
	}

	// Static files for frontend (production)
	// r.Static("/static", "./static")
	// r.StaticFile("/", "./static/index.html")

	return r
}
