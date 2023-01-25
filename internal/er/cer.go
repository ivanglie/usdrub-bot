package er

import (
	"fmt"
	"sort"
	"sync"

	br "github.com/ivanglie/go-br-client"
)

// CashExchangeRate exchange rate of cash.
type CashExchangeRate struct {
	sync.RWMutex
	branches     []br.Branch
	buyBranches  map[int][]string
	sellBranches map[int][]string
	buyMin       float64
	buyMax       float64
	buyAvg       float64
	sellMin      float64
	sellMax      float64
	sellAvg      float64
	err          error
	rateFunc     func() (*br.Rates, error)
}

func NewCashExchangeRate(rateFunc func() (*br.Rates, error)) *CashExchangeRate {
	return &CashExchangeRate{rateFunc: rateFunc}
}

// Update currency exchange cash rate.
func (c *CashExchangeRate) Update(wg *sync.WaitGroup) {
	update := func() {
		c.Lock()
		defer c.Unlock()

		var r *br.Rates
		r, c.err = c.rateFunc()
		c.branches = r.Branches
		c.buyMin, c.sellMin, c.buyMax, c.sellMax, c.buyAvg, c.sellAvg = findMma(c.branches)
		c.buyBranches, c.sellBranches = buyBranches(c.branches), sellBranches(c.branches)
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

// Rate of currency exchange cash returns of buyMin, buyMax, buyAvg, sellMin, sellMax, sellAvg.
func (c *CashExchangeRate) Rate() (float64, float64, float64, float64, float64, float64) {
	c.RLock()
	defer c.RUnlock()
	return c.buyMin, c.buyMax, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg
}

// String representation of currency exchange cash rate.
func (c *CashExchangeRate) String() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("Buy:\t%.2f .. %.2f RUB (avg %.2f)\nSell:\t%.2f .. %.2f RUB (avg %.2f)",
		c.buyMax, c.buyMin, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg)
}

// BuyBranches represented as string.
func (c *CashExchangeRate) BuyBranches() map[int][]string {
	c.RLock()
	defer c.RUnlock()
	return c.buyBranches
}

// SellBranches represented as string.
func (c *CashExchangeRate) SellBranches() map[int][]string {
	c.RLock()
	defer c.RUnlock()
	return c.sellBranches
}

func buyBranches(b []br.Branch) map[int][]string {
	sort.Sort(sort.Reverse(br.ByBuySorter(b)))
	d := []string{}
	i := 0
	for _, v := range b {
		if v.Buy != 0 {
			i++
			d = append(d, fmt.Sprintf("%d) %.2f RUB (_%v_): %s, %s, %s", i, v.Buy, v.Updated.Format("02.01.2006 15:04"), v.Bank, v.Address, v.Subway))
		}
	}
	return func(b []string, n int) map[int][]string {
		m := make(map[int][]string)
		j := 0
		for i := range b {
			if i%n == 0 {
				j = i + n

				var s []string
				if j < len(b) {
					s = b[i:j]
				} else {
					s = b[i:]
				}

				m[(j-n)/n] = s
			}
		}
		return m
	}(d, 5)
}

func sellBranches(b []br.Branch) map[int][]string {
	sort.Sort(br.BySellSorter(b))
	d := []string{}
	i := 0
	for _, v := range b {
		if v.Sell != 0 {
			i++
			d = append(d, fmt.Sprintf("%d) %.2f RUB (_%v_): %s, %s, %s", i, v.Sell, v.Updated.Format("02.01.2006 15:04"), v.Bank, v.Address, v.Subway))
		}
	}
	return func(b []string, n int) map[int][]string {
		m := make(map[int][]string)
		j := 0
		for i := range b {
			if i%n == 0 {
				j = i + n

				var s []string
				if j < len(b) {
					s = b[i:j]
				} else {
					s = b[i:]
				}

				m[(j-n)/n] = s
			}
		}
		return m
	}(d, 5)
}

// Find min, max and avg
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
