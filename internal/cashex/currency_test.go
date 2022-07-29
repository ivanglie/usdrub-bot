package cashex

import (
	"testing"
	"time"
)

func Test_min(t *testing.T) {
	type args struct {
		b []branch
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "min",
			args: args{[]branch{
				{"b", "a", "s", "c", 12.34, 58.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.34, 56.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.34, 57.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
			}},
			want: 56.78,
		},
		{
			name: "0",
			args: args{[]branch{}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.b); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_max(t *testing.T) {
	type args struct {
		b []branch
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "max",
			args: args{[]branch{
				{"b", "a", "s", "c", 12.34, 56.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.34, 58.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.34, 57.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
			}},
			want: 58.78,
		},
		{
			name: "0",
			args: args{[]branch{}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.b); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_avg(t *testing.T) {
	type args struct {
		b []branch
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "avg",
			args: args{[]branch{
				{"b", "a", "s", "c", 12.34, 30, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.34, 40, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
				{"b", "a", "s", "c", 12.34, 80, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()},
			}},
			want: 50,
		},
		{
			name: "0",
			args: args{[]branch{}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := avg(tt.args.b); got != tt.want {
				t.Errorf("avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataRace(t *testing.T) {
	c := New("%.2f", "moskva")
	go func() {
		for {
			c.Update()
		}
	}()

	for i := 0; i < 10; i++ {
		c.Format()
		time.Sleep(100 * time.Millisecond)
	}
}
