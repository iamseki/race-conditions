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
		// m.Lock() this can be done but not works at all, because this context has no iteraction with shared resource (variable total)
		go count(&wg, &m)
		// m.Unlock()
	}
	wg.Wait()

	return total
}

func countSameMemoryAddressVariableWithChannel() (total int) {
	wg := sync.WaitGroup{}
	c := make(chan bool, 1) // forces channel to have a buffer of one byte size, this will block other goroutines to writes it when buffer is already full

	count := func(wg *sync.WaitGroup, c chan bool) {
		c <- true // others go routine will stuck here untill the value is released by "<-c"
		total++
		<-c
		wg.Done()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go count(&wg, c)
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

	fmt.Println("******** concurrently count with channel ********")
	r = countSameMemoryAddressVariableWithChannel()
	fmt.Println("first call result with channel: ", r)

	r = countSameMemoryAddressVariableWithChannel()
	fmt.Println("second call result with channel: ", r)

	fmt.Println(`----------------------------------------
	now channel were using in concurrently count implementation, since the channel has a buffer size of 1 byte there is no way
	of goroutines writes on the same channel if it is full, so this will syncronize the execution.
	This is a hack !!!! The best use case for channel is for comunication between goroutines 1!! 
----------------------------------------`)
}
