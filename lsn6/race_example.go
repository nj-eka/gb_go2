package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	const rounds = 1000
	var (
		scores int32
		mu     sync.Mutex
		start  = make(chan struct{})
		wg     sync.WaitGroup
	)
	wg.Add(2 * rounds)
	for i := 0; i < rounds; i++ {
		go func() {
			defer wg.Done()
			<-start
			mu.Lock()
			defer mu.Unlock()
			scores++
		}()
		go func() {
			defer wg.Done()
			<-start
			atomic.AddInt32(&scores, -1)
		}()
	}
	close(start)
	wg.Wait()
	fmt.Println("final scores = ", scores)
}
