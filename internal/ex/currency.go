package ex

import (
	"fmt"
	"sync"
)

var lock sync.RWMutex

// Currency of exchange rate
type Currency struct {
	source   string
	format   string
	rate     float64
	err      error
	rateFunc func() (float64, error)
}

func NewCurrency(name, format string, rateFunc func() (float64, error)) (Currency, error) {
	lock.Lock()
	defer lock.Unlock()

	p := Currency{
		source:   name,
		format:   format,
		rateFunc: rateFunc,
	}
	p.rate, p.err = p.rateFunc()
	if p.err != nil {
		return p, fmt.Errorf("%s error: %v", p.source, p.err)
	}
	return p, nil
}

// Get formated exchange rate
func (p *Currency) Format() string {
	lock.RLock()
	defer lock.RUnlock()

	r := fmt.Sprintf(p.format, p.rate)
	if p.rate <= 0.0 {
		r = fmt.Sprintf("%s error: wrong value of rate=%.2f", p.source, p.rate)
	}
	return r
}
