package cash

import (
	"errors"
	"testing"
	"time"

	"github.com/ivanglie/usdrub-bot/pkg/bankiru-go"
	"github.com/stretchr/testify/assert"
)

func Test_rate_Update(t *testing.T) {
	r := Get()
	r.f = func() (*bankiru.Branches, error) {
		rates := &bankiru.Branches{
			Currency: "USD",
			City:     bankiru.Moscow,
			Items: []bankiru.Branch{
				{Bank: "b", Subway: "s", Currency: "c", Buy: 49.0, Sell: 51.0, Updated: time.Now()},
				{Bank: "b", Subway: "s", Currency: "c", Buy: 50.0, Sell: 52.0, Updated: time.Now()},
				{Bank: "b", Subway: "s", Currency: "c", Buy: 51.0, Sell: 53.0, Updated: time.Now()},
			},
		}

		return rates, nil
	}

	r.Update()
	assert.Equal(t, 3, len(r.branches))

	// Error
	r.f = func() (*bankiru.Branches, error) {
		rates := &bankiru.Branches{
			Currency: "USD",
			City:     bankiru.Moscow,
			Items: []bankiru.Branch{
				{Bank: "b", Subway: "s", Currency: "c", Buy: 49.0, Sell: 51.0, Updated: time.Now()},
				{Bank: "b", Subway: "s", Currency: "c", Buy: 50.0, Sell: 52.0, Updated: time.Now()},
			},
		}

		return rates, errors.New("error")
	}

	r.Update()
	assert.Equal(t, 3, len(r.branches))
}

func Test_rate_String(t *testing.T) {
	r := &cash{}
	r.branches = []bankiru.Branch{{Bank: "b", Subway: "s", Currency: "c", Buy: 100.0, Sell: 200.0, Updated: time.Now()}}

	assert.NotEmpty(t, r.String())
}

func Test_mma(t *testing.T) {
	// Min, max and avg
	b := []bankiru.Branch{
		{
			Bank:     "b",
			Subway:   "s",
			Currency: "c",
			Buy:      13.00,
			Sell:     58.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b",
			Subway:   "s",
			Currency: "c",
			Buy:      12.00,
			Sell:     56.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{
			Bank:     "b",
			Subway:   "s",
			Currency: "c",
			Buy:      14.00,
			Sell:     57.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
	}

	bmin, smin, bmax, smax, bavg, savg := mma(b)
	assert.Equal(t, bmin, 12.00)
	assert.Equal(t, smin, 56.00)
	assert.Equal(t, bmax, 14.00)
	assert.Equal(t, smax, 58.00)
	assert.Equal(t, bavg, 13.00)
	assert.Equal(t, savg, 57.00)

	// Empty branches
	b = []bankiru.Branch{}

	bmin, smin, bmax, smax, bavg, savg = mma(b)
	assert.Equal(t, bmin, 0.00)
	assert.Equal(t, smin, 0.00)
	assert.Equal(t, bmax, 0.00)
	assert.Equal(t, smax, 0.00)
	assert.Equal(t, bavg, 0.00)
	assert.Equal(t, savg, 0.00)
}

func Test_rate_BuyBranches(t *testing.T) {
	b := []bankiru.Branch{
		{
			Bank:     "b1",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b2",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b3",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b4",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b5",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b6",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
	}

	c := &cash{}
	c.buyBranches = buyBranches(b)
	bb := c.BuyBranches()

	assert.Equal(t, len(bb), 6)
}

func Test_rate_SellBranches(t *testing.T) {
	b := []bankiru.Branch{
		{
			Bank:     "b1",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b2",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b3",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b4",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b5",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
		{
			Bank:     "b6",
			Subway:   "s",
			Currency: "c",
			Buy:      1.5,
			Sell:     1.00,
			Updated:  func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }(),
		},
	}

	c := &cash{}
	c.sellBranches = sellBranches(b)
	sb := c.SellBranches()

	assert.Equal(t, len(sb), 6)
}
