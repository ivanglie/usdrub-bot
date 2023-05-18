package br

import (
	"encoding/json"
	"fmt"
	"time"
)

// Bank or branch.
type Branch struct {
	Bank     string    `json:"bank"`
	Address  string    `json:"address"`
	Subway   string    `json:"subway"`
	Currency string    `json:"currency"`
	Buy      float64   `json:"buy"`
	Sell     float64   `json:"sell"`
	Updated  time.Time `json:"updated"`
}

// Currency type.
type Currency string

// City type.
type City string

// Rates by banks and their branches.
type Rates struct {
	Currency Currency `json:"currency"`
	City     City     `json:"city"`
	Branches []Branch `json:"branches"`
}

// NewBranch creates a new Branch instance.
func newBranch(bank, address, subway, currency string, buy, sell float64, updated time.Time) Branch {
	return Branch{bank, address, subway, currency, buy, sell, updated}
}

// ByBuySorter implements sort.Interface based on the Buy field.
type ByBuySorter []Branch

// Len, Swap and Less implement sort.Interface for ByBuySorter.
func (b ByBuySorter) Len() int           { return len(b) }
func (b ByBuySorter) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByBuySorter) Less(i, j int) bool { return b[i].Buy < b[j].Buy }

// BySellSorter implements sort.Interface based on the Sell field.
type BySellSorter []Branch

// Len, Swap and Less implement sort.Interface for BySellSorter.
func (s BySellSorter) Len() int           { return len(s) }
func (s BySellSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySellSorter) Less(i, j int) bool { return s[i].Sell < s[j].Sell }

// String representation of cash currency exchange rates.
func (r *Rates) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(b)
}
