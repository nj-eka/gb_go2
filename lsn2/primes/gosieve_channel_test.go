package primes

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleSieveWithChannel() {
	fmt.Println(SieveWithChannel(100))
}

func BenchmarkSieveWithChannel(b *testing.B) {
	for i := 2; i < b.N; i++ {
		_ = SieveWithChannel(i)
	}
}

func TestSieveWithChannel(t *testing.T) {
	type args struct {
		top int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"2", 	args{2},[]int{2}},
		{"3", 	args{3},[]int{2,3}},
		{"4", 	args{4},[]int{2,3}},
		{"5", 	args{5},[]int{2,3,5}},
		{"10", 	args{10},[]int{2,3,5,7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SieveWithChannel(tt.args.top); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SieveWithBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
