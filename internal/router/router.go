package router

import (
	"aiksava-lb/internal/controllers"
	"aiksava-lb/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.SetTrustedProxies(nil)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	proxy := r.Group("/proxy")
	proxy.Any("/*proxyPath", controllers.ProxyHandler)

	return r
}