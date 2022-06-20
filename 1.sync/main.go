package main

import "sync"

// Problem: Race condition
// build: go build --race main.go

var (
	balance int = 100
)

func Deposit(amount int, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	defer lock.Unlock() // Unlock is a method of the Mutex struct

	lock.Lock() // Lock the mutex
	b := balance
	balance = b + amount
}

func Balance(lock *sync.RWMutex) int {
	lock.RLock() // Lock the mutex read from the mutex
	b := balance
	lock.RUnlock()
	return b
}

func main() {
	var wg sync.WaitGroup
	// var lock sync.Mutex // Mutex is a struct that implements the Lock and Unlock methods
	var lock sync.RWMutex // Mutex is a struct that implements the Lock and Unlock methods

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go Deposit(i*100, &wg, &lock)
	}

	wg.Wait()
	println(Balance(&lock))
}
