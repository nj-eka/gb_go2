//Написать программу, которая использует мьютекс для безопасного доступа к данным
//из нескольких потоков. Выполните трассировку программы
package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
)

type Counter struct {
	sync.Mutex
	value uint64
}

func (c *Counter) Increment() {
	c.Lock()
	defer c.Unlock()
	c.value++
}

func runMutexWorkers(countGoroutines int, countWorkers int) uint64 {
	var (
		wg    sync.WaitGroup
		wPool = make(chan struct{}, countWorkers)
		mc    = Counter{}
	)

	for i := 0; i < countGoroutines; i++ {
		wPool <- struct{}{}
		wg.Add(1)
		go func(id int) {
			defer func() {
				<-wPool
				wg.Done()
				fmt.Printf("Goroutine with id = %d - done\n", id)
			}()
			mc.Increment()
		}(i)
	}
	wg.Wait()
	fmt.Printf("mutex: Counter = %d\n", mc.value)
	return mc.value
}

func main() {
	const rounds, workers = 100000, 100
	if err := trace.Start(os.Stderr); err == nil {
		defer trace.Stop()
	}
	runMutexWorkers(rounds, workers)
}
