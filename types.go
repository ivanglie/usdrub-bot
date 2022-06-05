package main

import (
	"fmt"
	"sync"
)

// Provides code safe
var mu sync.RWMutex

// Source of exchange rate
type Source struct {
	name     string
	pattern  string
	rate     float64
	err      error
	rateFunc func() (float64, error)
}

// Update rate from source
func (s *Source) updateRate() error {
	mu.Lock()
	defer mu.Unlock()
	s.rate, s.err = s.rateFunc()
	if s.err != nil {
		return fmt.Errorf("%s error: %v", s.name, s.err)
	}
	return nil
}

// Get formated exchange rate from Source struct
func (s *Source) getRatef() string {
	mu.RLock()
	defer mu.RUnlock()
	res := fmt.Sprintf(s.pattern, s.rate)
	if s.err != nil || s.rate <= 0.0 {
		res = fmt.Sprintf("%s error: %v", s.name, s.err)
	}
	return res
}
