package moex

import (
	"testing"
	"time"
)

func TestNewCurrency(t *testing.T) {
	type args struct {
		name    string
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Equals",
			args: args{name: "any", pattern: "%.2f"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.pattern)
			got.Update()
			if got.rate <= 0 {
				t.Errorf("NewCurrency() = %v", got)
			}
		})
	}
}

func Test_dataRace(t *testing.T) {
	c := New("%.2f")
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
