package main

import (
	"testing"
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
	if got := messageByCommand(1, "start"); got.ChatID != 1 && len(got.Text) == 0 {
		t.Errorf("messageByCommand() = %v", got)
	}

	if got := messageByCommand(1, "forex"); got.ChatID != 1 && len(got.Text) == 0 {
		t.Errorf("messageByCommand() = %v", got)
	}

	if got := messageByCommand(1, "moex"); got.ChatID != 1 && len(got.Text) == 0 {
		t.Errorf("messageByCommand() = %v", got)
	}

	if got := messageByCommand(1, "cbrf"); got.ChatID != 1 && len(got.Text) == 0 {
		t.Errorf("messageByCommand() = %v", got)
	}

	if got := messageByCommand(1, "cash"); got.ChatID != 1 && len(got.Text) == 0 {
		t.Errorf("messageByCommand() = %v", got)
	}

	if got := messageByCommand(1, "help"); got.ChatID != 1 && len(got.Text) == 0 {
		t.Errorf("messageByCommand() = %v", got)
	}

	if got := messageByCommand(1, " "); got.ChatID != 1 && len(got.Text) == 0 {
		t.Errorf("messageByCommand() = %v", got)
	}
}
