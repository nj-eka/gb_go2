// С помощью пула воркеров написать программу, которая запускает 1000 горутин,
// каждая из которых увеличивает число на 1.
// Дождаться завершения всех горутин и убедиться, что при каждом запуске программы итоговое число равно 1000.

package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value uint64
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
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
