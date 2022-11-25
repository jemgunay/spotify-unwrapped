package auth

import (
	"sync"
	"time"
)

// Access is a concurrency safe wrapper around the Spotify access token.
type Access struct {
	token   string
	refresh refreshFunc
	expiry  time.Time
	mu      *sync.RWMutex
}

type refreshFunc func() (string, time.Time, error)

// New initialises a new Access.
func New(refreshFunc refreshFunc) *Access {
	return &Access{
		refresh: refreshFunc,
		mu:      &sync.RWMutex{},
	}
}

// Get retrieves the access token, or lazy-fetches a fresh one if it has expired.
func (a *Access) Get() (string, error) {
	a.mu.RLock()
	currentExpiry := a.expiry
	a.mu.RUnlock()

	now := time.Now().UTC()
	var err error
	if now.After(currentExpiry) {
		err = a.Refresh()
	}

	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.token, err
}

// Refresh refreshes the access token.
func (a *Access) Refresh() error {
	token, expiry, err := a.refresh()
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.token = token
	a.expiry = expiry
	return nil
}
