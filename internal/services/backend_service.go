package services

import (
	"aiksava-lb/internal/config"
	"aiksava-lb/internal/models"
	"fmt"
	"log"
	"os"
)

var Backends []*models.Backend
var roundRobin models.RoundRobin = models.RoundRobin{
	Count: 0,
}

func InitBackends() {
	fmt.Println("Initializing backends...", os.Getenv("ENVIRONMENT"), config.PRODUCTION)

	allBackends := []string{
		"https://aiksava.onrender.com",
		"https://aiksava-s1.onrender.com",
		"https://aiksava-s2.onrender.com",
	}

	if os.Getenv("ENVIRONMENT") != config.PRODUCTION {
		allBackends = append(
			[]string{"http://localhost:8081"},
			allBackends...,
		)
	}

	for _, url := range allBackends {
		Backends = append(Backends, &models.Backend{
			URL: url,
		})
	}
}

func PickBackend() (*models.Backend, int) {
	LogUserCountsPerBackend()

	maxBackendCount := config.GetBackendCount()

	limit := min(maxBackendCount, len(Backends))

	for i := range limit {
		backend := Backends[i]

		if backend.ActiveUsers < config.GetMaxUsersPerServer() {
			fmt.Println("Assigning to backend:", backend.URL, "Current users:", backend.ActiveUsers)
			backend.ActiveUsers++
			return backend, (i + 1) % limit
		}
	}

	resBackend := Backends[roundRobin.Count]
	Backends[roundRobin.Count].ActiveUsers++
	fmt.Println("All backends at capacity. Assigning to backend (round-robin):", resBackend.URL, roundRobin.Count)
	roundRobin.Count = (roundRobin.Count + 1) % limit
	return resBackend, roundRobin.Count
}

func DecrementBackend(url string) {
	for _, b := range Backends {
		if b.URL == url && b.ActiveUsers > 0 {
			b.ActiveUsers--
		}
	}
}

func LogUserCountsPerBackend() {
	log.Println("User counts per backend:")
	for _, backend := range Backends {
		log.Printf("  %s: %d users", backend.URL, backend.ActiveUsers)
	}
}
