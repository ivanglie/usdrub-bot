package cexrate

import (
	"fmt"
	"sort"
	"sync"
	"time"

	br "github.com/ivanglie/go-br-client"
)

const (
	Prefix = "Exchange rates of cash"
	Suffix = "in branches in Moscow, Russia by Banki.ru"
)

// rate represents currency exchange rate of cash.
type rate struct {
	sync.RWMutex
	name         string
	f            func() (*br.Rates, error)
	branches     []br.Branch
	buyBranches  []string
	sellBranches []string
	buyMin       float64
	buyMax       float64
	buyAvg       float64
	sellMin      float64
	sellMax      float64
	sellAvg      float64
	err          error
	errDate      time.Time
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
		RateInstance = &rate{name: Prefix, f: func() (*br.Rates, error) { return br.NewClient().Rates(br.USD, br.Moscow) }}
	}

	return RateInstance
}

// Update exchange rate of cash.
func (r *rate) Update() {
	r.Lock()
	defer r.Unlock()

	rates, err := r.f()
	if rates == nil || err != nil {
		r.err = err
		r.errDate = time.Now()
		return
	}

	r.err = nil
	r.branches = rates.Branches
	r.buyMin, r.sellMin, r.buyMax, r.sellMax, r.buyAvg, r.sellAvg = findMma(r.branches)
	r.buyBranches, r.sellBranches = buyBranches(r.branches), sellBranches(r.branches)
}

// String representation of currency exchange cash rate.
func (r *rate) String() string {
	r.RLock()
	defer r.RUnlock()

	return fmt.Sprintf("Buy:\t%.2f .. %.2f RUB (avg %.2f)\nSell:\t%.2f .. %.2f RUB (avg %.2f)",
		r.buyMax, r.buyMin, r.buyAvg, r.sellMin, r.sellMax, r.sellAvg)
}

// BuyBranches represented as string.
func (r *rate) BuyBranches() []string {
	r.RLock()
	defer r.RUnlock()

	return r.buyBranches
}

// SellBranches represented as string.
func (r *rate) SellBranches() []string {
	r.RLock()
	defer r.RUnlock()

	return r.sellBranches
}

// buyBranches represented as string.
func buyBranches(b []br.Branch) []string {
	sort.Sort(sort.Reverse(br.ByBuySorter(b)))

	s := []string{}
	for _, v := range b {
		if v.Buy != 0 {
			s = append(s, fmt.Sprintf("%.2f RUB (%v): %s, %s, %s", v.Buy, v.Updated.Format("02.01.2006 15:04"), v.Bank, v.Address, v.Subway))
		}
	}

	return s
}

// sellBranches represented as string.
func sellBranches(b []br.Branch) []string {
	sort.Sort(br.BySellSorter(b))

	s := []string{}
	for _, v := range b {
		if v.Sell != 0 {
			s = append(s, fmt.Sprintf("%.2f RUB (%v): %s, %s, %s", v.Sell, v.Updated.Format("02.01.2006 15:04"), v.Bank, v.Address, v.Subway))
		}
	}

	return s
}

// findMma returns min, max and average values of buy and sell rates.
func findMma(b []br.Branch) (bmin, smin, bmax, smax, bavg, savg float64) {
	if len(b) == 0 {
		return
	}

	btotal, stotal := float64(0), float64(0)

	bb, sb := []br.Branch{}, []br.Branch{}
	for _, v := range b {
		if v.Buy != 0 {
			bb = append(bb, v)
		}

		if v.Sell != 0 {
			sb = append(sb, v)
		}
	}

	bmin, bmax = bb[0].Buy, bb[0].Buy
	for _, v := range bb {
		if v.Buy < bmin {
			bmin = v.Buy
		}
		if v.Buy > bmax {
			bmax = v.Buy
		}
		btotal += v.Buy
	}

	smin, smax = sb[0].Sell, sb[0].Sell
	for _, v := range sb {
		if v.Sell < smin {
			smin = v.Sell
		}
		if v.Sell > smax {
			smax = v.Sell
		}
		stotal += v.Sell
	}

	bavg, savg = btotal/float64(len(bb)), stotal/float64(len(sb))

	return
}
