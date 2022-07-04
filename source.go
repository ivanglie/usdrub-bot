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
	r, e := s.rateFunc()
	s.setRate(r)
	s.setErr(e)
	if e != nil {
		return fmt.Errorf("%s error: %v", s.getName(), s.getErr())
	}
	return nil
}

// Get formated exchange rate from Source struct
func (s *Source) getRatef() string {
	res := fmt.Sprintf(s.getPattern(), s.getRate())
	if s.getErr() != nil || s.getRate() <= 0.0 {
		res = fmt.Sprintf("%s error: %v", s.getName(), s.getErr())
	}
	return res
}

func (s *Source) getName() string {
	mu.RLock()
	defer mu.RUnlock()
	return s.name
}

func (s *Source) getPattern() string {
	mu.RLock()
	defer mu.RUnlock()
	return s.pattern
}

func (s *Source) getRate() float64 {
	mu.RLock()
	defer mu.RUnlock()
	return s.rate
}

func (s *Source) getErr() error {
	mu.RLock()
	defer mu.RUnlock()
	return s.err
}

func (s *Source) setRate(rate float64) {
	mu.Lock()
	defer mu.Unlock()
	s.rate = rate
}

func (s *Source) setErr(err error) {
	mu.Lock()
	defer mu.Unlock()
	s.err = err
}
