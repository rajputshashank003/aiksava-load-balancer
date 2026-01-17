package controllers

import (
	"aiksava-lb/internal/config"
	"aiksava-lb/internal/services"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var proxyTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 30 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	MaxIdleConns:          100,
	MaxIdleConnsPerHost:   10,
	IdleConnTimeout:       90 * time.Second,
}

func ProxyHandler(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	fmt.Println("sessions counts:")
	fmt.Println(services.GetSessionCountsPerBackend())

	if c.Request.Method == "OPTIONS" {
		setCORSHeaders(c.Writer, origin)
		c.AbortWithStatus(204)
		return
	}

	sessionID := c.GetHeader(config.SessionHeader)
	session, exist, actualSessionID := services.GetOrCreateSession(sessionID)

	if !exist {
		fmt.Print("Creating new session with ID: ", actualSessionID, "\n")
		backend, nextInt := services.PickBackend()
		services.ColdStartAtInd(nextInt)
		session = services.CreateSession(actualSessionID, backend.URL)
	} else {
		fmt.Println("exist krta hai ")
	}

	services.TouchSession(session)
	backendURL := session.Backend

	fmt.Printf("Proxying to backend: %s\n", backendURL)

	target, err := url.Parse(backendURL)
	if err != nil {
		setCORSHeaders(c.Writer, origin)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid backend URL",
		})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = proxyTransport

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		req.Header.Set("Host", target.Host)
		req.Header.Set("X-Forwarded-Host", c.Request.Host)
		req.Header.Set("X-Forwarded-Proto", "https")
		req.Header.Set("X-Forwarded-For", c.ClientIP())
		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", "Aiksava-LB")
		}
		req.Header.Set(config.SessionHeader, session.ID)
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		headersToRemove := []string{}
		for key := range resp.Header {
			lowerKey := strings.ToLower(key)
			if strings.HasPrefix(lowerKey, "access-control") {
				headersToRemove = append(headersToRemove, key)
			}
		}
		for _, header := range headersToRemove {
			resp.Header.Del(header)
		}

		if origin != "" {
			resp.Header.Set("Access-Control-Allow-Origin", origin)
		} else {
			resp.Header.Set("Access-Control-Allow-Origin", "*")
		}
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
		resp.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Forwarded-For, X-Requested-With, access_token, "+config.SessionHeader)
		resp.Header.Set("Access-Control-Expose-Headers", config.SessionHeader)
		resp.Header.Set("Vary", "Origin")

		resp.Header.Set(config.SessionHeader, actualSessionID)

		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("Proxy error: %v\n", err)
		setCORSHeaders(w, origin)
		w.Header().Set(config.SessionHeader, session.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf(`{"error": "Backend unavailable: %s"}`, err.Error())))
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func setCORSHeaders(w http.ResponseWriter, origin string) {
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Forwarded-For, X-Requested-With, access_token, "+config.SessionHeader)
	w.Header().Set("Access-Control-Expose-Headers", config.SessionHeader)
	w.Header().Set("Vary", "Origin")
}
