package router

import (
	"aiksava-lb/internal/config"
	"aiksava-lb/internal/controllers"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())
	r.SetTrustedProxies([]string{"0.0.0.0/0"})

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if config.IsAllowedOrigin(origin) {
				return true
			}
			return origin == "http://localhost:8000" ||
				origin == "http://127.0.0.1:8000" ||
				origin == "http://[::1]:8000" ||
				origin == "https://aiksava.onrender.com" ||
				origin == "https://aiksava.vercel.app"
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"access_token",
			config.SessionHeader,
		},
		ExposeHeaders:    []string{"Content-Length", config.SessionHeader},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Proxy
	api := r.Group("/api")
	api.Any("/*apiPath", controllers.ProxyHandler)

	return r
}
