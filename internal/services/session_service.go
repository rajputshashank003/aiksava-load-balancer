package services

import (
	"aiksava-lb/internal/models"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	Sessions = make(map[string]*models.Session)
	Mu       sync.Mutex
)

func GetOrCreateSession(sessionID string) (*models.Session, bool, string) {
	Mu.Lock()
	defer Mu.Unlock()
	fmt.Println("GetOrCreateSession called with sessionID:", sessionID)
	session, exists := Sessions[sessionID]

	if exists {
		return session, exists, sessionID
	}
	sessionID = GenerateSessionID()
	return session, exists, sessionID
}

func CreateSession(sessionID string, backendURL string) *models.Session {
	Mu.Lock()
	defer Mu.Unlock()
	fmt.Println("Creating session with ID:", sessionID, "for backend:", backendURL)
	session := &models.Session{
		ID:       sessionID,
		Backend:  backendURL,
		LastSeen: time.Now(),
	}

	Sessions[sessionID] = session
	return session
}

func TouchSession(session *models.Session) {
	Mu.Lock()
	defer Mu.Unlock()

	session.LastSeen = time.Now()
}

func GenerateSessionID() string {
	newSessionId := uuid.NewString()
	return newSessionId
}

func GetSessionCountsPerBackend() map[string]int {
	Mu.Lock()
	defer Mu.Unlock()

	counts := make(map[string]int)
	for _, session := range Sessions {
		counts[session.Backend]++
	}

	for backend, count := range counts {
		fmt.Printf("%s : %d count \n", backend, count)
	}

	return counts
}
