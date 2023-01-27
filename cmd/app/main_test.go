package main

import (
	"reflect"
	"testing"
	"unsafe"

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

	t.Log("cash.BuyBranches()", cash.BuyBranches())
	t.Log("cash.SellBranches()=", cash.SellBranches())

	for _, v := range []string{"Buy", "BuyMore", "Sell", "SellMore", "Help", " "} {
		if got := messageByCallbackData(1, v); got.ChatID != 1 && len(got.Text) == 0 {
			t.Errorf("messageByCallbackData() = %v", got)
		}
	}
}
