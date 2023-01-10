package cashex

import (
	"os"
	"path/filepath"
	"strings"
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
		name  string
		args  args
		want  float64
		want1 float64
		want2 float64
		want3 float64
		want4 float64
		want5 float64
	}{
		{
			name: "Min, max and avg",
			args: args{[]branch{
				{"b", "a", "s", "c", 13.00, 58.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.00, 56.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 14.00, 57.00, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
			}},
			want:  12,
			want1: 56,
			want2: 14,
			want3: 58,
			want4: 13,
			want5: 57,
		},
		{
			name:  "Empty branches",
			args:  args{[]branch{}},
			want:  0,
			want1: 0,
			want2: 0,
			want3: 0,
			want4: 0,
			want5: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, got4, got5 := findMma(tt.args.b)
			if got != tt.want {
				t.Errorf("mma() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("mma() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("mma() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("mma() got3 = %v, want %v", got3, tt.want3)
			}
			if got4 != tt.want4 {
				t.Errorf("mma() got4 = %v, want %v", got4, tt.want4)
			}
			if got5 != tt.want5 {
				t.Errorf("mma() got5 = %v, want %v", got5, tt.want5)
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
		t.Errorf("b is nil")
	}

	branchesCount := len(currency.branches)
	buyBranchesCount := len(strings.Split(buyBranches(currency.branches), "\n"))
	sellBranchesCount := len(strings.Split(sellBranches(currency.branches), "\n"))

	if branchesCount != 1 {
		t.Errorf("branchesCount got = %v, want %v", branchesCount, 3)
	}

	if buyBranchesCount != 1 {
		t.Errorf("buyBranchesCount got = %v, want %v", buyBranchesCount, 2)
	}

	if sellBranchesCount != 1 {
		t.Errorf("sellBranchesCount got = %v, want %v", sellBranchesCount, 2)
	}
}
