package er

import (
	"fmt"
	"sync"
)

// ExchangeRate of exchange rate
type ExchangeRate struct {
	sync.RWMutex
	rate     float64
	err      error
	rateFunc func() (float64, error)
}

func NewExchangeRate(rateFunc func() (float64, error)) *ExchangeRate {
	return &ExchangeRate{rateFunc: rateFunc}
}

// Update
func (c *ExchangeRate) Update(wg *sync.WaitGroup) {
	update := func() {
		c.Lock()
		defer c.Unlock()
		c.rate, c.err = c.rateFunc()
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
func (c *ExchangeRate) Rate() (float64, error) {
	c.RLock()
	defer c.RUnlock()
	return c.rate, c.err
}

// Get formatted exchange rate
func (c *ExchangeRate) String() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("%.2f RUB", c.rate)
}
