package cashex

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_dataRace(t *testing.T) {
	c := New("moskva")
	go func() {
		for {
			c.Update(nil)
		}
	}()

	for i := 0; i < 10; i++ {
		c.Rate()
		time.Sleep(100 * time.Millisecond)
	}
}

func Test_mma(t *testing.T) {
	type args struct {
		b []branch
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
			args: args{[]branch{
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
			args: args{[]branch{
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

func Test_parseBranches(t *testing.T) {
	currency := Currency{}

	dir, _ := os.Getwd()
	absFilePath := filepath.Join(dir, "../../test/bankiru")

	currency.parseBranches("file:" + absFilePath)

	if len(currency.branches) == 0 {
		t.Errorf("currency.branches is empty")
	}

	branchesCount := len(currency.branches)

	buyBranches := buyBranches(currency.branches)
	buyBranchesCount := 0
	for _, v := range buyBranches {
		buyBranchesCount = buyBranchesCount + len(v)
	}

	sellBranches := sellBranches(currency.branches)
	sellBranchesCount := 0
	for _, v := range sellBranches {
		sellBranchesCount = sellBranchesCount + len(v)
	}

	if branchesCount != 5 {
		t.Errorf("branchesCount got = %v, want %v", branchesCount, 5)
	}

	if buyBranchesCount != 4 {
		t.Errorf("buyBranchesCount got = %v, want %v", buyBranchesCount, 4)
	}

	if sellBranchesCount != 4 {
		t.Errorf("sellBranchesCount got = %v, want %v", sellBranchesCount, 4)
	}
}
