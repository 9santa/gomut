package main

import (
	"fmt"
	"sync"
	"time"
)

/* ===== Benchmark =====
time=3.291264104s, counter=320000

3.95user 0.11system 0:03.35elapsed 120%CPU (0avgtext+0avgdata 58848maxresident)k
0inputs+3336outputs (0major+4267minor)pagefaults 0swaps

Go's builtin mutexes are still a bit faster than Futex version
*/

func main() {
	goroutines := 8
	workTime := time.Duration(10) * time.Microsecond
	iterations := 40000

	mu := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(goroutines)
	start := make(chan struct{}) // gate

	counter := 0

	begin := time.Now()
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			<-start // every goroutine blocks here
			for j := 0; j < iterations; j++ {
				mu.Lock()
				counter++
				simulateWork(workTime)
				mu.Unlock()
			}
		}()
	}

	close(start) // goroutines unblock and run
	wg.Wait()
	elapsed := time.Since(begin)
	fmt.Printf("time=%v, counter=%d\n", elapsed, counter)
}

func simulateWork(d time.Duration) {
	if d <= 0 {
		return
	}
	endTime := time.Now().Add(d)
	for time.Now().Before(endTime) {
		// simulate
	}
}
