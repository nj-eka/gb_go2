package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	//	runtime.GOMAXPROCS(1)

	const (
		timeout = 1 * time.Second
		cars    = 10
	)

	if err := trace.Start(os.Stderr); err == nil {
		defer trace.Stop()
	}

	ctx, stop := context.WithTimeout(context.Background(), timeout)
	defer stop()

	type result struct {
		carNumber, scores int
	}
	start := make(chan struct{})
	finish := make(chan result, cars)

	var wg sync.WaitGroup
	wg.Add(cars)

	for i := 0; i < cars; i++ {
		go func(num int) {
			defer wg.Done()
			<-start
			scores := 0
			for {
				select {
				case <-ctx.Done():
					finish <- result{num, scores}
					return
				default:
					scores++
					if num%2 == 0 && scores%100 == 0 {
						runtime.Gosched()
					}
				}
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(finish)
	}()
	close(start)
	for res := range finish {
		fmt.Printf("#%d: %d\n", res.carNumber, res.scores)
	}
}
