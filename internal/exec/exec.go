package exec

import "sync"

var globalMutex = sync.Mutex{}

func Safe(f func()) {
	globalMutex.Lock()
	f()
	globalMutex.Unlock()
}

func ParallelN(count int, f func(n int)) {
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(n int) {
			f(n)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func Parallel(a, b func()) {
	ParallelN(2, func(n int) {
		if n == 0 {
			a()
		} else {
			b()
		}
	})
}
