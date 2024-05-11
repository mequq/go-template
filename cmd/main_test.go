package main

import (
	"log/slog"
	"reflect"
	"testing"
)

func Test_initSlogLogger(t *testing.T) {
	type args struct {
		cfg LogingConfig
	}
	tests := []struct {
		name string
		args args
		want *slog.Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initSlogLogger(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initSlogLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
