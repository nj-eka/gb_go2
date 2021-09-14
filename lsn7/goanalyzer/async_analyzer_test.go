package goanalyzer

import (
	"testing"
)

func TestGetGoFuncCallCount(t *testing.T) {
	type args struct {
		fileName string
		funcName string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"GetGoFuncCallCount4sieve_channel",
		args{"../../lsn2/primes/gosieve_channel.go", "filter"},
		1,
		false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := GetGoFuncCallCount(tt.args.fileName, tt.args.funcName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGoFuncCallCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("GetGoFuncCallCount() gotCount = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}
