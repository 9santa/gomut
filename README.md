===== BENCHMARKS =====
- Spinlock
time=3.249180753s, counter=320000

16.93user 1.60system 0:03.27elapsed 566%CPU (0avgtext+0avgdata 21028maxresident)k
0inputs+32outputs (0major+1926minor)pagefaults 0swaps

- Futexlock
time=3.44880997s, counter=320000

4.27user 0.31system 0:03.47elapsed 131%CPU (0avgtext+0avgdata 21256maxresident)k
0inputs+0outputs (0major+1877minor)pagefaults 0swaps

As we can see, user&system time dropped massively.
This is due to futexWait 'parking' the thread, while parked they consume ~0 CPU.
While Spinlock version keeps doing atomic.SwapUint32(...) in a tight loop, a lot of CPU work.

- Go's builtin
time=3.291264104s, counter=320000

3.95user 0.11system 0:03.35elapsed 120%CPU (0avgtext+0avgdata 58848maxresident)k
0inputs+3336outputs (0major+4267minor)pagefaults 0swaps

Go's builtin mutexes are still a bit faster than Futex version
