package concurent

import (
	"fmt"
	"runtime"
	"sync"
)

var DefGoroutines = runtime.NumCPU()

func InPlaceMap[T any](goroutines int, data []T, function func(*T)) error {
	if goroutines == 0 {
		return fmt.Errorf("goroutines = 0")
	}
	if len(data) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(start, step int, datSize int) {
			defer wg.Done()

			for index := start; index < datSize; index += step {
				function(&data[index])
			}
		}(i, goroutines, len(data))
	}

	wg.Wait()
	return nil
}

func DefInPlaceMap[T any](data []T, function func(*T)) {
	err := InPlaceMap(DefGoroutines, data, function)
	if err != nil {
		panic(err)
	}
}

func Map[S any, R any](goroutines int, data []S, function func(S) R) ([]R, error) {
	if goroutines == 0 {
		return nil, fmt.Errorf("goroutines = 0")
	}
	if len(data) == 0 {
		return nil, nil
	}

	result := make([]R, len(data))
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(start, step int, datSize int) {
			defer wg.Done()

			for index := start; index < datSize; index += step {
				result[index] = function(data[index])
			}
		}(i, goroutines, len(data))
	}

	wg.Wait()
	return result, nil
}

func DefMap[S any, R any](data []S, function func(S) R) ([]R, error) {
	result, err := Map(DefGoroutines, data, function)
	if err != nil {
		panic(err)
	}
	return result, nil
}
