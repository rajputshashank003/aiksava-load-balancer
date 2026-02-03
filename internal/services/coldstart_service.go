package services

import (
	"aiksava-lb/internal/models"
	"log"
	"net/http"
	"time"
)

func fetchBackend(backend *models.Backend) {
	client := &http.Client{
		Timeout: 360 * time.Second,
	}

	resp, err := client.Get(backend.URL + "/health")
	if err != nil {
		log.Println("Backend unreachable:", backend.URL, err)
		return
	}
	defer resp.Body.Close()

	log.Println("Backend warmed:", backend.URL, resp.Status)
}

func ColdStart() {
	for _, backend := range Backends {
		go fetchBackend(backend)
	}
}

func FirstColdStart() {
	go fetchBackend(Backends[0])
}

func ColdStartAtInd(ind int) {
	go fetchBackend(Backends[ind])
}

func HealthCheck() {
	// todo: uncomment this to start cron job for health checks
	// ticker := time.NewTicker(12 * time.Minute)

	// go func() {
	// 	for range ticker.C {
	// 		fetchBackend(Backends[0])
	// 	}
	// }()
}