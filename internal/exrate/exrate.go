package exrate

import (
	"fmt"
	"sync"
	"time"

	"github.com/ivanglie/go-cbr-client"
	"github.com/ivanglie/go-coingate-client"
	"github.com/ivanglie/go-moex-client"
)

const (
	Prefix = "1 US Dollar equals"

	Forex = "Forex"
	MOEX  = "Moscow Exchange"
	CBRF  = "Russian Central Bank"
)

// rate represents exchange rate.
type rate struct {
	sync.RWMutex
	name    string
	f       func() (float64, error)
	value   float64
	err     error
	errDate time.Time
}

// update exchange rate.
func (r *rate) update() {
	r.Lock()
	defer r.Unlock()

	v, err := r.f()
	if err != nil || v == 0 {
		r.err = err
		r.errDate = time.Now()
		return
	}

	r.value = v
	r.err = nil
}

// String representation of rate.
func (r *rate) String() string {
	r.RLock()
	defer r.RUnlock()

	return fmt.Sprintf("%.2f RUB by %s", r.value, r.name)
}

// rates represents exchange rates.
type rates struct {
	sync.RWMutex
	values []*rate
}

var (
	ratesInstance *rates
	lock          = &sync.Mutex{}
)

// Get returns instance of Rates.
func Get() *rates {
	lock.Lock()
	defer lock.Unlock()

	if ratesInstance == nil {
		ratesInstance = &rates{}
		ratesInstance.values = []*rate{
			{name: Forex, f: func() (float64, error) { return coingate.NewClient().GetRate("USD", "RUB") }},
			{name: MOEX, f: func() (float64, error) { return moex.NewClient().GetRate(moex.USDRUB) }},
			{name: CBRF, f: func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) }}}
	}

	return ratesInstance
}

// Update exchange rates.
func (r *rates) Update() {
	r.Lock()
	defer r.Unlock()

	for _, v := range r.values {
		v.update()
	}
}

// Value returns rate by name.
func (r *rates) Value(name string) *rate {
	r.RLock()
	defer r.RUnlock()

	for _, value := range r.values {
		if value.name == name {
			return value
		}
	}

	return nil
}

// String representation of rates.
func (r *rates) String() string {
	r.RLock()
	defer r.RUnlock()

	var s string
	for _, v := range r.values {
		s += fmt.Sprintf("%s\n", v)
	}

	return s
}
