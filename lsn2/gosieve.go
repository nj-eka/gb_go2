// 'primes' package contains various functions for working (search, definition, etc.) with prime numbers.
//
// 1. 'Sieve' function prints all prime numbers in range from 1 up to specified number (using goroutines).
package primes

import "fmt"

// Send the sequence 2, 3, 4, … , 'top' to channel 'ch'.
func generate(top int, ch chan<- int) {
	for i := 2; i <= top; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'src' to channel 'dst',
// removing those divisible by 'prime'.
func filter(src <-chan int, dst chan<- int, prime int) {
	for i := range src { // Loop over values received from 'src'.
		if i%prime != 0 {
			dst <- i // Send 'i' to channel 'dst'.
		}
	}
}

// The prime sieve: Daisy-chain filter processes together.
func Sieve(top int) {
	ch := make(chan int) // Create a new channel.
	go generate(top, ch) // Start generate() as a subprocess.
	for {
		prime := <-ch
		fmt.Print(prime, "\n")
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}
