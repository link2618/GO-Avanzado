package main

import (
	"fmt"
	"sync"
	"time"
)

func ExpensiveFibonacci(n interface{}) int {
	switch n := n.(type) {
	case int:
		if n <= 1 {
			return n
		}
		return ExpensiveFibonacci(n-1) + ExpensiveFibonacci(n-2)
	default:
		return 0
	}
}

// Memory holds a function and a map of results
type Memory struct {
	f     Function                       // Function to be used
	cache map[interface{}]FunctionResult // Map of results for a given key
	lock  sync.RWMutex                   // Lock to protect the cache

	InProgress map[interface{}]bool                  // Map of jobs in progress
	IsPending  map[interface{}][]chan FunctionResult // Map of jobs waiting for a response
}

// A function has to recive a value and return a value and an error
type Function func(key interface{}) (interface{}, error)

// The result of a function
type FunctionResult struct {
	value interface{}
	err   error
}

// NewCache creates a new cache
func NewCache(f Function) *Memory {
	return &Memory{
		f:          f,
		cache:      make(map[interface{}]FunctionResult),
		lock:       sync.RWMutex{},
		InProgress: make(map[interface{}]bool),
		IsPending:  make(map[interface{}][]chan FunctionResult),
	}
}

func (m *Memory) service(key interface{}) (interface{}, error) {
	// Check if the job is already in progress
	m.lock.RLock()
	_, ok := m.InProgress[key]
	m.lock.RUnlock()

	// If the job is already in progress, then wait for the response
	if ok {
		// If the job is already in progress, then wait for the response
		response := make(chan FunctionResult)
		defer close(response)

		// Add the channel to the pending list
		m.lock.Lock()
		m.IsPending[key] = append(m.IsPending[key], response)
		m.lock.Unlock()

		// Wait for the response
		res := <-response
		return res.value, res.err
	}

	// If the job is not in progress, then start the job
	m.lock.Lock()
	m.InProgress[key] = true
	m.lock.Unlock()

	// Start the job
	fmt.Printf("Starting job %d\n", key)
	fnresult, err := m.f(key)
	res := FunctionResult{value: fnresult, err: err}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Get the pending workers for this job
	m.lock.RLock()
	pendingWorkers, exist := m.IsPending[key]
	m.lock.RUnlock()

	// If there are pending workers, then send the response
	if exist {
		for _, worker := range pendingWorkers {
			worker <- res
		}
	}

	// We are done with this job, reset the state
	m.lock.Lock()
	m.InProgress[key] = false
	m.IsPending[key] = make([]chan FunctionResult, 0)
	m.lock.Unlock()

	// Add the value to the cache
	m.lock.Lock()
	m.cache[key] = res
	m.lock.Unlock()

	fmt.Printf("Finished job %d, got %d\n", key, res)
	return res.value, res.err

}

// Get returns the value for a given key
func (m *Memory) Get(key interface{}) (interface{}, error) {

	// Lock the cache
	m.lock.RLock()

	// Check if the value is in the cache
	res, exist := m.cache[key]

	// Unlock the cache
	m.lock.RUnlock()

	// If the value is in the cache, return it
	if exist {
		return res.value, res.err
	}

	// If the value is not in the cache, then start the service
	res.value, res.err = m.service(key)
	return res.value, res.err

}

// Function to be used in the cache
func GetFibonacci(key interface{}) (interface{}, error) {
	return ExpensiveFibonacci(key), nil
}

func main() {
	empezar := time.Now()
	// Create a cache and some values
	cache := NewCache(GetFibonacci)
	values := []int{46, 46, 42, 42, 41, 41, 46, 46, 46, 42, 42, 46}

	var wg sync.WaitGroup

	// For each value to calculate, get the value and print the time it took to calculate
	for _, v := range values {
		go func(v int) {
			defer wg.Done()

			start := time.Now()

			res, err := cache.Get(v)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			fmt.Printf("%v:%d took %v\n", v, res, time.Since(start))
		}(v)
		wg.Add(1)
	}
	wg.Wait()

	// This is to prove that cache is working
	fmt.Println()
	fmt.Println("Doing it all again!")
	for _, v := range values {
		go func(v int) {
			defer wg.Done()

			start := time.Now()

			res, err := cache.Get(v)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			fmt.Printf("%v:%d took %v\n", v, res, time.Since(start))
		}(v)
		wg.Add(1)
	}

	wg.Wait()

	fmt.Printf("Total time: %v\n", time.Since(empezar))
}
