package internal

import "sync"

var globalMutex sync.Mutex

func ExecSafe(f func()) {
	globalMutex.Lock()
	f()
	globalMutex.Unlock()
}

func ExecParallel(count int, f func(index int)) {
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(index int) {
			f(index)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
