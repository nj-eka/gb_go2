// С помощью пула воркеров написать программу, которая запускает 1000 горутин,
// каждая из которых увеличивает число на 1.
// Дождаться завершения всех горутин и убедиться, что при каждом запуске программы итоговое число равно 1000.

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func runAtomicWorkers(countGoroutines int, countWorkers int) (counter uint64) {
	var (
		wg    sync.WaitGroup
		wPool = make(chan struct{}, countWorkers)
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
			atomic.AddUint64(&counter, 1)
		}(i)
	}
	wg.Wait()
	fmt.Printf("atomic: Counter = %d\n", counter)
	return
}
