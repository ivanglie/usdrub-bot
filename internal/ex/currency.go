package ex

import (
	"fmt"
	"sync"
)

// Currency of exchange rate
type Currency struct {
	sync.RWMutex
	pattern  string
	rate     float64
	err      error
	rateFunc func() (float64, error)
}

func New(pattern string, rateFunc func() (float64, error)) *Currency {
	return &Currency{
		pattern:  pattern,
		rateFunc: rateFunc,
	}
}

func (c *Currency) Update() {
	c.Lock()
	defer c.Unlock()
	c.rate, c.err = c.rateFunc()
}

// Get formated exchange rate
func (c *Currency) Format() string {
	c.RLock()
	defer c.RUnlock()
	r := fmt.Sprintf(c.pattern, c.rate)
	if c.rate <= 0.0 {
		r = fmt.Sprintf("ex error: wrong value of rate=%.2f", c.rate)
	}
	return r
}
