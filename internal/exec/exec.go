package exec

import "sync"

var globalMutex = sync.Mutex{}

// Safe executes the given function inside a global mutex.
// Deadlock will occur if "f" also calls "Safe".
func Safe(f func()) {
	globalMutex.Lock()
	f()
	globalMutex.Unlock()
}

// ParallelN is a concurrency helper for iterating through
// any typed slice in parallel.
func ParallelN(n int, f func(n int)) {
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(n int) {
			f(n)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// Parallel executes all input functions in parallel.
func Parallel(fns ...func()) {
	ParallelN(len(fns), func(n int) {
		fns[n]()
	})
}
