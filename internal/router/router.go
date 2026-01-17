package router

import (
	"aiksava-lb/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// Trust all proxies for now (you can restrict this in production)
	r.SetTrustedProxies(nil)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS is handled in the proxy controller to avoid duplicate headers
	
	r.Any("/*proxyPath", controllers.ProxyHandler)

	return r
}
