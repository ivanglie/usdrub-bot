package cryptoexrate

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ivanglie/usdrub-bot/pkg/go-bestchange-client"
)

const (
	Prefix = "1 USDT equals"
	Suffix = "in Moscow, Russia by bestchange.com"
)

// rate represents currency exchange rate of cash.
type rate struct {
	sync.RWMutex
	name    string
	f       func() (float64, error)
	value   float64
	err     error
	errDate time.Time
}

var (
	RateInstance *rate
	lock         = &sync.Mutex{}
)

// Get returns instance of Rate.
func Get() *rate {
	lock.Lock()
	defer lock.Unlock()

	if RateInstance == nil {
		RateInstance = &rate{name: Prefix, f: func() (float64, error) { return bestchange.NewClient().Rate(bestchange.Moscow) }}
	}

	return RateInstance
}

// Update exchange rate of cash.
func (r *rate) Update() {
	r.Lock()
	defer r.Unlock()

	v, err := r.f()
	if v == 0 || err != nil {
		log.Printf("[ERROR] %s: value=%v, error=%v", r.name, v, err)

		r.err = err
		r.errDate = time.Now()
		return
	}

	r.err = nil
	r.value = v
}

// String representation of currency exchange cash rate.
func (r *rate) String() string {
	r.RLock()
	defer r.RUnlock()

	return fmt.Sprintf("%.2f RUB %s", r.value, Suffix)
}
