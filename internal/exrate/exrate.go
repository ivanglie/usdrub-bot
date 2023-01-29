package exrate

import (
	"fmt"
	"log"
	"sync"
)

// Rate of exchange rate
type Rate struct {
	sync.RWMutex
	rate     float64
	err      error
	rateFunc func() (float64, error)
}

func NewRate(rateFunc func() (float64, error)) *Rate {
	return &Rate{rateFunc: rateFunc}
}

// Update
func (c *Rate) Update(wg *sync.WaitGroup) {
	update := func() {
		c.Lock()
		defer c.Unlock()
		c.rate, c.err = c.rateFunc()
		if c.err != nil {
			log.Println(c.err)
		}
	}

	if wg == nil {
		update()
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		update()
	}()
}

// Get exchange rate
func (c *Rate) Rate() (float64, error) {
	c.RLock()
	defer c.RUnlock()
	return c.rate, c.err
}

// Get formatted exchange rate
func (c *Rate) String() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("%.2f RUB", c.rate)
}
