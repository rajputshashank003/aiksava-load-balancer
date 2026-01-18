package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {

	allowedOrigins := map[string]bool{
		"https://aiksava.vercel.app":  true,
	}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		fmt.Println("Request Origin:", origin)
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Authorization, X-Session-Id",
		)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}