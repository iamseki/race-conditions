package main

import (
	"fmt"
	"sync"
)

func countSameMemoryAddressVariable() (total int) {
	wg := sync.WaitGroup{}

	count := func(wg *sync.WaitGroup) {
		total++
		wg.Done()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go count(&wg)
	}
	wg.Wait()

	return total
}

func countSameMemoryAddressVariableWithMutex() (total int) {
	wg := sync.WaitGroup{}
	m := sync.Mutex{}

	count := func(wg *sync.WaitGroup, m *sync.Mutex) {
		m.Lock()
		total++
		m.Unlock()
		wg.Done()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go count(&wg, &m)
	}
	wg.Wait()

	return total
}

func main() {
	fmt.Println("******** simple concurrently count ********")
	r := countSameMemoryAddressVariable()
	fmt.Println("first call result: ", r)

	r = countSameMemoryAddressVariable()
	fmt.Println("second call result: ", r)

	fmt.Println(`----------------------------------------
	doesnt matters how many times I call the func, even if the max iteration is 1000 
	and the expected result seems to be 1000 since there is no synchronization between the goroutines, 
	the result can always be different.
----------------------------------------`)

	fmt.Println("******** concurrently count with mutex ********")
	r = countSameMemoryAddressVariableWithMutex()
	fmt.Println("first call result with mutex: ", r)

	r = countSameMemoryAddressVariableWithMutex()
	fmt.Println("second call result with mutex: ", r)

	fmt.Println(`----------------------------------------
	now mutex were using in concurrently count implementation, there is no way of goroutines access the same address at the same time
	so the result is always 1000 as expected.
----------------------------------------`)
}
