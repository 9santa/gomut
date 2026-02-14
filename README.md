## Benchmarks

All runs were executed with:

- `goroutines * iterations = counter`
- Timing printed by the program is **wall-clock time** for the critical section loop.
- `/usr/bin/time` output reflects the whole `go run ...` command (compile + link + run).

---

### Spinlock
time=3.249180753s, counter=320000

16.93user 1.60system 0:03.27elapsed 566%CPU (0avgtext+0avgdata 21028maxresident)k
0inputs+32outputs (0major+1926minor)pagefaults 0swaps

### Futex lock
time=3.44880997s, counter=320000

4.27user 0.31system 0:03.47elapsed 131%CPU (0avgtext+0avgdata 21256maxresident)k
0inputs+0outputs (0major+1877minor)pagefaults 0swaps

**Observation:** `user` + `system` time drops massively compared to the spinlock.

**Why:** `futexWait` *parks* the thread when the lock is contended. While parked, it consumes ~0 CPU.  
In contrast, the spinlock version keeps retrying `atomic.SwapUint32(...)` in a tight loop, burning CPU cycles.

---

### Go’s built-in `sync.Mutex`
time=3.291264104s, counter=320000

3.95user 0.11system 0:03.35elapsed 120%CPU (0avgtext+0avgdata 58848maxresident)k
0inputs+3336outputs (0major+4267minor)pagefaults 0swaps


**Result:** Go’s built-in mutex is still a bit faster than the futex-based lock in this benchmark.
