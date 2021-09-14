package primes

// Send the sequence 2, 3, 4, … [top] to channel 'ch'
func generate(top int, ch chan<- int) {
	for i := 2; i <= top; i++ {
		ch <- i // Send 'i' to channel 'ch'
	}
	close(ch)
}

// filter - copy values from channel 'src' to channel 'dst'
// removing those divisible by 'prime'
func filter(src <-chan int, dst chan<- int, prime int) {
	for i := range src {  // Loop over values received from 'src'
		if i % prime != 0 {
			dst <- i  // Send 'i' to channel 'dst'
		}
	}
	close(dst)
}

// SieveWithChannel - finds prime numbers in range from 2 up to [top] using go concurrency with daisy-chain filter pattern - the slowest function !!!
// BenchmarkSieveWithChannel-4         1000           1212408 ns/op           13345 B/op         94 allocs/op
func SieveWithChannel(top int) []int {
	primes := make([]int, 0, top)
	ch := make(chan int)
	go generate(top, ch)
	for {
		if prime, more := <-ch; more {
			primes = append(primes, prime)
			ch1 := make(chan int)
			go filter(ch, ch1, prime)
			ch = ch1
		} else { break }
	}
	return primes
}
