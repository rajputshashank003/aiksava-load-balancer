package services

import (
	"aiksava-lb/internal/models"
	"log"
	"net/http"
	"time"
)

func fetchBackend(backend *models.Backend) {
	client := &http.Client{
		Timeout: 90 * time.Second,
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

func ColdStartAtInd(ind int) {
	go fetchBackend(Backends[ind])
}
