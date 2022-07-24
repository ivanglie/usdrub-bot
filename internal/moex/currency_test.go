package moex

import (
	"testing"
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
			if got := NewCurrency(tt.args.name, tt.args.pattern); got.rate <= 0 {
				t.Errorf("NewCurrency() = %v", got)
			}
		})
	}
}
