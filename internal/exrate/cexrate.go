package exrate

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	br "github.com/ivanglie/go-br-client"
)

// CashRate exchange rate of cash.
type CashRate struct {
	sync.RWMutex
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
}

// New exchange rate of cash.
func NewCashRate(rates *br.Rates, err error) *CashRate {
	if rates == nil || err != nil {
		return &CashRate{}
	}

	r := &CashRate{}
	r.Lock()
	defer r.Unlock()
	r.err = err
	r.branches = rates.Branches
	r.buyMin, r.sellMin, r.buyMax, r.sellMax, r.buyAvg, r.sellAvg = findMma(r.branches)
	r.buyBranches, r.sellBranches = buyBranches(r.branches), sellBranches(r.branches)

	return r
}

// Update exchange rate of cash.
func UpdateCashRate(f func() (*br.Rates, error)) *CashRate {
	r, err := f()
	return NewCashRate(r, err)
}

// Rate of currency exchange cash returns of buyMin, buyMax, buyAvg, sellMin, sellMax, sellAvg.
func (c *CashRate) Rate() (float64, float64, float64, float64, float64, float64, error) {
	c.RLock()
	defer c.RUnlock()
	if c.buyMin == 0 && c.buyMax == 0 && c.buyAvg == 0 && c.sellMin == 0 && c.sellMax == 0 && c.sellAvg == 0 {
		return 0, 0, 0, 0, 0, 0, errors.New("values of c.buyMin, c.buyMax, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg are 0")
	}
	return c.buyMin, c.buyMax, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg, nil
}

// String representation of currency exchange cash rate.
func (c *CashRate) String() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("Buy:\t%.2f .. %.2f RUB (avg %.2f)\nSell:\t%.2f .. %.2f RUB (avg %.2f)",
		c.buyMax, c.buyMin, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg)
}

// BuyBranches represented as string.
func (c *CashRate) BuyBranches() []string {
	c.RLock()
	defer c.RUnlock()
	return c.buyBranches
}

// SellBranches represented as string.
func (c *CashRate) SellBranches() []string {
	c.RLock()
	defer c.RUnlock()
	return c.sellBranches
}

func buyBranches(branch []br.Branch) []string {
	sort.Sort(sort.Reverse(br.ByBuySorter(branch)))

	s := []string{}
	for _, b := range branch {
		if b.Buy != 0 {
			s = append(s, fmt.Sprintf("%.2f RUB (%v): %s, %s, %s", b.Buy, b.Updated.Format("02.01.2006 15:04"), b.Bank, b.Address, b.Subway))
		}
	}

	return s
}

func sellBranches(branch []br.Branch) []string {
	sort.Sort(br.BySellSorter(branch))

	s := []string{}
	for _, v := range branch {
		if v.Sell != 0 {
			s = append(s, fmt.Sprintf("%.2f RUB (%v): %s, %s, %s", v.Sell, v.Updated.Format("02.01.2006 15:04"), v.Bank, v.Address, v.Subway))
		}
	}

	return s
}

// Find min, max and avg.
func findMma(r []br.Branch) (bmin, smin, bmax, smax, bavg, savg float64) {
	if len(r) == 0 {
		return
	}

	btotal, stotal := float64(0), float64(0)

	bb, sb := []br.Branch{}, []br.Branch{}
	for _, v := range r {
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
