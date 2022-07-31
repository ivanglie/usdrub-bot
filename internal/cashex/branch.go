package cashex

import (
	"time"
)

// branch
type branch struct {
	Bank     string    `json:"bank"`
	Address  string    `json:"address"`
	Subway   string    `json:"subway"`
	Currency string    `json:"currency"`
	Buy      float64   `json:"buy"`
	Sell     float64   `json:"sell"`
	Updated  time.Time `json:"updated"`
}

func newBranch(bank, address, subway, currency string, buy, sell float64, updated time.Time) branch {
	return branch{
		Bank:     bank,
		Address:  address,
		Subway:   subway,
		Currency: currency,
		Buy:      buy,
		Sell:     sell,
		Updated:  updated}
}

// ByBuySorter implements sort.Interface based on the Buy field
type ByBuySorter []branch

func (b ByBuySorter) Len() int           { return len(b) }
func (b ByBuySorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByBuySorter) Less(i, j int) bool { return b[i].Buy < b[j].Buy }

// BySellSorter implements sort.Interface based on the Sell field
type BySellSorter []branch

func (s BySellSorter) Len() int           { return len(s) }
func (s BySellSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySellSorter) Less(i, j int) bool { return s[i].Sell < s[j].Sell }
