package main

import (
	"sync"
)

// Branch
type Branch struct {
	mu      sync.RWMutex
	Bank    string  `json:"bank"`
	Address string  `json:"address"`
	Subway  string  `json:"subway"`
	Code    string  `json:"code"`
	Buy     float64 `json:"buy"`
	Sell    float64 `json:"sell"`
	Updated string  `json:"updated"`
}

func (b *Branch) getSell() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Sell
}
