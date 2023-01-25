package er

import (
	"testing"
	"time"

	br "github.com/ivanglie/go-br-client"
)

func Test_mma(t *testing.T) {
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
