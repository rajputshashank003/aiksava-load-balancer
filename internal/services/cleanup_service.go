package services

import (
	"aiksava-lb/internal/config"
	"fmt"
	"time"
)

func StartSessionCleanup() {
	ttl := config.GetSessionTTL()

	go func() {
		for {
			time.Sleep(30 * time.Second)

			Mu.Lock()
			now := time.Now()

			for id, session := range Sessions {
				if now.Sub(session.LastSeen) > ttl {
					fmt.Println("Deleting session due to TTL expiry:", id)
					DecrementBackend(session.Backend)
					delete(Sessions, id)
				}
			}

			Mu.Unlock()
		}
	}()
}
