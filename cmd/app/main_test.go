package main

import "testing"

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
