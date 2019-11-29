package internal

import "sync"

var globalMutex = sync.Mutex{}

func ExecSafe(f func()) {
	globalMutex.Lock()
	f()
	globalMutex.Unlock()
}

func ExecParallel(count int, f func(n int)) {
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
