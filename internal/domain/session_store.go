package domain

import (
	"sync"
	"time"
)

var (
	sessionStore     SessionStore
	sessionStoreOnce sync.Once
)

// GetSessionStore returns the singleton instance of the session store
func GetSessionStore() SessionStore {
	sessionStoreOnce.Do(func() {
		sessionStore = newInMemorySessionStore()
	})
	return sessionStore
}

// inMemorySessionStore implements SessionStore using an in-memory map
type inMemorySessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func newInMemorySessionStore() *inMemorySessionStore {
	return &inMemorySessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *inMemorySessionStore) Get(id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[id]
	if !exists {
		return nil, nil
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		delete(s.sessions, id)
		return nil, nil
	}

	return session, nil
}

func (s *inMemorySessionStore) Set(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.ID] = session
	return nil
}

func (s *inMemorySessionStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)
	return nil
}
