package primes

import (
	"fmt"
	"github.com/yourbasic/bit"
	"reflect"
	"testing"
)

func ExampleSieveWithBits() {
	fmt.Println(SieveWithBits(10).String())
}

func BenchmarkSieveWithBits(b *testing.B) {
	for i := 2; i < b.N; i++ {
		_ = SieveWithBits(i)
	}
}

func TestSieveWithBits(t *testing.T) {
	type args struct {
		top int
	}
	tests := []struct {
		name string
		args args
		want *bit.Set
	}{
		{"2", 	args{2}, bit.New(2)},
		{"3", 	args{3}, bit.New(2, 3)},
		{"4", 	args{4},bit.New(2,3)},
		{"5", 	args{5},bit.New(2,3,5)},
		{"10", 	args{10},bit.New(2,3,5,7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SieveWithBits(tt.args.top); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SieveWithBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
