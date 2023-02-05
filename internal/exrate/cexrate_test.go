package exrate

import (
	"testing"
	"time"

	br "github.com/ivanglie/go-br-client"
)

func TestNewCashRate(t *testing.T) {
	if got := NewCashRate(&br.Rates{}, nil); got == nil {
		t.Errorf("NewCashRate() = %v", got)
	}
}

func TestCashRate_String(t *testing.T) {
	r := &CashRate{}
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

func TestCashRate_BuyBranches(t *testing.T) {
	b := []br.Branch{
		{"b1", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b2", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b3", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b4", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b5", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b6", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
	}

	c := &CashRate{}
	c.buyBranches = buyBranches(b)
	bb := c.BuyBranches()

	if got := len(bb); got != 2 {
		t.Errorf("len(c.BuyBranches()) = %v, want %v", got, 2)
	}

	if got := len(bb[0]); got != 5 {
		t.Errorf("len(c.BuyBranches()[0]) = %v, want %v", got, 5)
	}

	if got := len(bb[1]); got != 1 {
		t.Errorf("len(c.BuyBranches()[1]) = %v, want %v", got, 1)
	}
}

func TestCashRate_SellBranches(t *testing.T) {
	b := []br.Branch{
		{"b1", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b2", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b3", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b4", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b5", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b6", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
		{"b7", "a", "s", "c", 1.5, 1.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
	}

	c := &CashRate{}
	c.sellBranches = sellBranches(b)
	sb := c.SellBranches()

	if got := len(sb); got != 2 {
		t.Errorf("len(c.SellBranches()) = %v, want %v", got, 2)
	}

	if got := len(sb[0]); got != 5 {
		t.Errorf("len(c.SellBranches()[0]) = %v, want %v", got, 5)
	}

	if got := len(sb[1]); got != 2 {
		t.Errorf("len(c.SellBranches()[1]) = %v, want %v", got, 2)
	}
}

func TestCashRate_Rate(t *testing.T) {
	c := &CashRate{buyMin: 0, buyMax: 0, buyAvg: 0, sellMin: 0, sellMax: 0, sellAvg: 0}
	if _, _, _, _, _, _, err := c.Rate(); err == nil {
		t.Errorf("CashRate.Rate() error = %v", err)
	}

	c = &CashRate{buyMin: 0, buyMax: 0, buyAvg: 0, sellMin: 0, sellMax: 1, sellAvg: 0}
	if _, _, _, _, _, _, err := c.Rate(); err != nil {
		t.Errorf("CashRate.Rate() error = %v", err)
	}
}
