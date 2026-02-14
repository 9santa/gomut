//go:build linux

package gomut

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Constants were taken from https://github.com/torvalds/linux/blob/master/include/uapi/linux/futex.h#L11
const (
	FUTEX_WAIT         = 0
	FUTEX_WAKE         = 1
	FUTEX_PRIVATE_FLAG = 128
)

func futexWait(addr *uint32, val uint32) error {
	_, _, err := unix.Syscall6(unix.SYS_FUTEX,
		uintptr(unsafe.Pointer(addr)),
		uintptr(FUTEX_WAIT|FUTEX_PRIVATE_FLAG),
		uintptr(val),
		0, 0, 0,
	)
	return err
}

func futexWake(addr *uint32, val uint32) error {
	_, _, err := unix.Syscall6(unix.SYS_FUTEX,
		uintptr(unsafe.Pointer(addr)),
		uintptr(FUTEX_WAKE|FUTEX_PRIVATE_FLAG),
		uintptr(val),
		0, 0, 0,
	)
	return err
}

type FutexLock struct {
	state uint32
}

func (l *FutexLock) Lock() {
	for {
		for atomic.LoadUint32(&l.state) == 1 {
			futexWait(&l.state, 1)
		}
		if atomic.SwapUint32(&l.state, 1) == 0 {
			return
		}
	}
}

func (l *FutexLock) Unlock() {
	atomic.StoreUint32(&l.state, 0)
	futexWake(&l.state, 1)
}

/* ===== Benchmark =====
time=3.44880997s, counter=320000

4.27user 0.31system 0:03.47elapsed 131%CPU (0avgtext+0avgdata 21256maxresident)k
0inputs+0outputs (0major+1877minor)pagefaults 0swaps

As we can see, user&system time dropped massively.
This is due to futexWait 'parking' the thread, while parked they consume ~0 CPU.
While Spinlock version keeps doing atomic.SwapUint32(...) in a tight loop, a lot of CPU work.
*/

func main() {
	goroutines := 32
	workTime := time.Duration(10) * time.Microsecond
	iterations := 10000

	lock := &FutexLock{}
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
