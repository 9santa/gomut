package gomut

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Spinlock struct {
	state uint32 // 0=unlocked, 1=locked
}

func (l *Spinlock) Lock() {
	for atomic.SwapUint32(&l.state, 1) == 1 {
		runtime.Gosched() // whenever a thread can't get a lock, it yields to another goroutine
		// wait... lock is already acquired by someone
	}
}

func (l *Spinlock) Unlock() {
	atomic.StoreUint32(&l.state, 0)
}

/* ===== Benchmark =====
time=3.249180753s, counter=320000

16.93user 1.60system 0:03.27elapsed 566%CPU (0avgtext+0avgdata 21028maxresident)k
0inputs+32outputs (0major+1926minor)pagefaults 0swaps
*/

func main() {
	goroutines := 8
	workTime := time.Duration(10) * time.Microsecond
	iterations := 40000

	lock := &Spinlock{}
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
				lock.Lock()
				counter++
				simulateWork(workTime)
				lock.Unlock()
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
