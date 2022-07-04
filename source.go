package main

import (
	"fmt"
	"sync"
)

// Source of exchange rate
type Source struct {
	mu       sync.RWMutex
	name     string
	pattern  string
	rate     float64
	err      error
	rateFunc func() (float64, error)
}

// Update rate from source
func (s *Source) UpdateRate() error {
	r, e := s.rateFunc()
	s.SetRate(r)
	s.SetErr(e)
	if e != nil {
		return fmt.Errorf("%s error: %v", s.GetName(), s.GetErr())
	}
	return nil
}

// Get formated exchange rate from Source struct
func (s *Source) GetRatef() string {
	res := fmt.Sprintf(s.GetPattern(), s.GetRate())
	if s.GetErr() != nil || s.GetRate() <= 0.0 {
		res = fmt.Sprintf("%s error: %v", s.GetName(), s.GetErr())
	}
	return res
}

func (s *Source) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.name
}

func (s *Source) GetPattern() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pattern
}

func (s *Source) GetRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rate
}

func (s *Source) GetErr() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.err
}

func (s *Source) SetRate(rate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rate = rate
}

func (s *Source) SetErr(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = err
}
