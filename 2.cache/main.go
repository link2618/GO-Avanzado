package main

import (
	"fmt"
	"sync"
	"time"
)

// Function to calculate fibonacci
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// Memory holds a function and a map of results
type Memory struct {
	f     Function               // Function to be used
	cache map[int]FunctionResult // Map of results for a given key
	lock  sync.RWMutex           // Lock to protect the cache
}

// A function has to recive a value and return a value and an error
type Function func(key int) (interface{}, error)

// The result of a function
type FunctionResult struct {
	value interface{}
	err   error
}

// NewCache creates a new cache
func NewCache(f Function) *Memory {
	return &Memory{f, make(map[int]FunctionResult), sync.RWMutex{}}
}

// Get returns the value for a given key
func (m *Memory) Get(key int) (interface{}, error) {

	// Lock the cache
	m.lock.Lock()

	// Check if the value is in the cache
	res, exist := m.cache[key]

	// Unlock the cache
	m.lock.Unlock()

	// If the value is not in the cache, calculate it
	if !exist {
		m.lock.Lock()
		res.value, res.err = m.f(key) // Calculate the value
		m.cache[key] = res            // Store the value in the cache
		m.lock.Unlock()
	}

	return res.value, res.err
}

// Function to be used in the cache
func GetFibonacci(key int) (interface{}, error) {
	return Fibonacci(key), nil
}

func main() {
	empezar := time.Now()
	// Create a cache and some values
	cache := NewCache(GetFibonacci)
	values := []int{42, 40, 41, 42, 38, 41, 42}

	var wg sync.WaitGroup

	maxGoroutines := 2
	channel := make(chan int, maxGoroutines)

	// For each value to calculate, get the value and print the time it took to calculate
	for _, v := range values {

		go func(index int) {
			defer wg.Done()
			channel <- 1
			start := time.Now()

			res, err := cache.Get(index)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			fmt.Printf("%v:%d took %v\n", index, res, time.Since(start))
			<-channel
		}(v)
		wg.Add(1)
	}

	wg.Wait()
	fmt.Printf("Tiempo total: %v\n", time.Since(empezar))

}
