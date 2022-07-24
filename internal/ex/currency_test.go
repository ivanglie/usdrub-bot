package ex

import (
	"reflect"
	"testing"
)

func TestNewCurrency(t *testing.T) {
	type args struct {
		name     string
		format   string
		rateFunc func() (float64, error)
	}
	tests := []struct {
		name    string
		args    args
		want    Currency
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCurrency(tt.args.name, tt.args.format, tt.args.rateFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCurrency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}
