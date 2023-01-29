package main

import (
	"errors"
	"math/rand"
	"reflect"
	"testing"
	"unsafe"

	br "github.com/ivanglie/go-br-client"
	"github.com/ivanglie/usdrub-bot/internal/exrate"
)

func Test_setupLog(t *testing.T) {
	type args struct {
		dbg bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ErrorLevel",
			args: args{dbg: false},
		},
		{
			name: "DebugLevel",
			args: args{dbg: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupLog(tt.args.dbg)
		})
	}
}

func Test_messageByCommand(t *testing.T) {
	for _, v := range []string{"start", "forex", "moex", "cbrf", "cash", "help", " "} {
		if got := messageByCommand(1, v); got.ChatID != 1 && len(got.Text) == 0 {
			t.Errorf("messageByCommand() = %v", got)
		}
	}
}

func Test_messageByCallbackData(t *testing.T) {
	cash = &exrate.CashRate{}
	v := reflect.ValueOf(cash)
	val := reflect.Indirect(v)

	bb := val.FieldByName("buyBranches")
	ptrToBb := unsafe.Pointer(bb.UnsafeAddr())
	realPtrToBb := (*map[int][]string)(ptrToBb)
	*realPtrToBb = map[int][]string{0: {"1", "2", "3", "4", "5"}, 1: {"6", "7", "8", "9", "10"}, 2: {"0"}}

	sb := val.FieldByName("sellBranches")
	ptrToSb := unsafe.Pointer(sb.UnsafeAddr())
	realPtrToSb := (*map[int][]string)(ptrToSb)
	*realPtrToSb = map[int][]string{0: {"10", "9", "8", "7", "6"}, 1: {"5", "4", "3", "2", "1"}, 2: {"0"}}

	for _, v := range []string{"Buy", "BuyMore", "Sell", "SellMore", "Help", " "} {
		if got := messageByCallbackData(1, v); got.ChatID != 1 && len(got.Text) == 0 {
			t.Errorf("messageByCallbackData() = %v", got)
		}
	}
}

func Test_updateRates(t *testing.T) {
	setupLog(false)

	mx = exrate.NewRate(func() (float64, error) { return 100 * rand.Float64(), nil })
	fx = exrate.NewRate(func() (float64, error) { return 100 * rand.Float64(), nil })
	cbrf = exrate.NewRate(func() (float64, error) { return 100 * rand.Float64(), nil })
	cash = exrate.NewCashRate(func() (*br.Rates, error) { return &br.Rates{}, nil })

	updateRates()

	mx = exrate.NewRate(func() (float64, error) { return 0, errors.New("error") })
	fx = exrate.NewRate(func() (float64, error) { return 0, errors.New("error") })
	cbrf = exrate.NewRate(func() (float64, error) { return 0, errors.New("error") })
	cash = exrate.NewCashRate(func() (*br.Rates, error) { return &br.Rates{}, errors.New("error") })

	updateRates()
}
