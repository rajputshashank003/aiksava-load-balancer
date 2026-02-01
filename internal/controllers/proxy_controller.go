package controllers

import (
	"aiksava-lb/internal/config"
	"aiksava-lb/internal/services"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(c *gin.Context) {
	sessionID := c.GetHeader(config.SessionHeader)
	session, exists, actualSessionID := services.GetOrCreateSession(sessionID)

	if !exists {
		backend, idx := services.PickBackend()
		services.ColdStartAtInd(idx)
		session = services.CreateSession(actualSessionID, backend.URL)
	}
	services.TouchSession(session)

	target, err := url.Parse(session.Backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid backend"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.Header.Set(config.SessionHeader, session.ID)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		for key := range resp.Header {
			if strings.HasPrefix(strings.ToLower(key), "access-control") {
				resp.Header.Del(key)
			}
		}
		resp.Header.Set(config.SessionHeader, actualSessionID)
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"backend unavailable"}`))
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
