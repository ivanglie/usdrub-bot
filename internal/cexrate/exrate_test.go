package cexrate

import (
	"errors"
	"testing"
	"time"

	"github.com/ivanglie/go-br-client"
	"github.com/stretchr/testify/assert"
)

func Test_rate_Update(t *testing.T) {
	r := Get()
	r.f = func() (*br.Rates, error) {
		rates := &br.Rates{
			Currency: br.USD,
			City:     br.Moscow,
			Branches: []br.Branch{
				{Bank: "b", Address: "a", Subway: "s", Currency: "c", Buy: 49.0, Sell: 51.0, Updated: time.Now()},
				{Bank: "b", Address: "a", Subway: "s", Currency: "c", Buy: 50.0, Sell: 52.0, Updated: time.Now()},
				{Bank: "b", Address: "a", Subway: "s", Currency: "c", Buy: 51.0, Sell: 53.0, Updated: time.Now()},
			},
		}

		return rates, nil
	}

	r.Update()
	assert.Equal(t, 3, len(r.branches))

	// Error
	r.f = func() (*br.Rates, error) {
		rates := &br.Rates{
			Currency: br.USD,
			City:     br.Moscow,
			Branches: []br.Branch{
				{Bank: "b", Address: "a", Subway: "s", Currency: "c", Buy: 49.0, Sell: 51.0, Updated: time.Now()},
				{Bank: "b", Address: "a", Subway: "s", Currency: "c", Buy: 50.0, Sell: 52.0, Updated: time.Now()},
			},
		}

		return rates, errors.New("error")
	}

	r.Update()
	assert.Equal(t, 3, len(r.branches))
}

func Test_rate_String(t *testing.T) {
	r := &rate{}
	r.branches = []br.Branch{{Bank: "b", Address: "a", Subway: "s", Currency: "c", Buy: 100.0, Sell: 200.0, Updated: time.Now()}}

	if got := r.String(); len(got) == 0 {
		t.Errorf("CashRate.String() = %v", got)
	}
}

func Test_findMma(t *testing.T) {
	type args struct {
		b []br.Branch
	}
	tests := []struct {
		name     string
		args     args
		bminWant float64
		sminWant float64
		bmaxWant float64
		smaxWant float64
		bavgWant float64
		savgWant float64
	}{
		{
			name: "Min, max and avg",
			args: args{[]br.Branch{
				{"b", "a", "s", "c", 13.00, 58.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.00, 56.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 14.00, 57.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
			}},
			bminWant: 12,
			sminWant: 56,
			bmaxWant: 14,
			smaxWant: 58,
			bavgWant: 13,
			savgWant: 57,
		},
		{
			name: "Min, max and avg without zeros values",
			args: args{[]br.Branch{
				{"b", "a", "s", "c", 13.00, 00.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 00.00, 00.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 14.00, 57.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
			}},
			bminWant: 13,
			sminWant: 57,
			bmaxWant: 14,
			smaxWant: 57,
			bavgWant: 13.5,
			savgWant: 57,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, got4, got5 := findMma(tt.args.b)
			if got != tt.bminWant {
				t.Errorf("mma() got = %v, want %v", got, tt.bminWant)
			}
			if got1 != tt.sminWant {
				t.Errorf("mma() got1 = %v, want %v", got1, tt.sminWant)
			}
			if got2 != tt.bmaxWant {
				t.Errorf("mma() got2 = %v, want %v", got2, tt.bmaxWant)
			}
			if got3 != tt.smaxWant {
				t.Errorf("mma() got3 = %v, want %v", got3, tt.smaxWant)
			}
			if got4 != tt.bavgWant {
				t.Errorf("mma() got4 = %v, want %v", got4, tt.bavgWant)
			}
			if got5 != tt.savgWant {
				t.Errorf("mma() got5 = %v, want %v", got5, tt.savgWant)
			}
		})
	}
}

func Test_rate_BuyBranches(t *testing.T) {
	b := []br.Branch{
		{"b1", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b2", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b3", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b4", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b5", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b6", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
	}

	c := &rate{}
	c.buyBranches = buyBranches(b)
	bb := c.BuyBranches()

	if got := len(bb); got != 6 {
		t.Errorf("len(c.BuyBranches()) = %v, want %v", got, 2)
	}
}

func Test_rate_SellBranches(t *testing.T) {
	b := []br.Branch{
		{"b1", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b2", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b3", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b4", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b5", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b6", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b7", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
	}

	c := &rate{}
	c.sellBranches = sellBranches(b)
	sb := c.SellBranches()

	if got := len(sb); got != 7 {
		t.Errorf("len(c.SellBranches()) = %v, want %v", got, 2)
	}
}
