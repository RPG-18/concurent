package concurrent

import (
	"fmt"
	"runtime"
	"sync"
)

// ReduceOption reduce order
type ReduceOption int

const (
	UnorderedReduce ReduceOption = iota + 1 // UnorderedReduce unordered
	OrderedReduce                           // OrderedReduce ordered
)

// DefGoroutines default goroutines count(runtime.NumCPU())
var DefGoroutines = runtime.NumCPU()

// Map calls function once for each item in sequence. The function takes a reference to the item, so that any modifications done to the item will appear in sequence.
func Map[T any](goroutines int, data []T, function func(*T)) error {
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

// DefMap calls function once for each item in sequence. The function takes a reference to the item, so that any modifications done to the item will appear in sequence.
func DefMap[T any](data []T, function func(*T)) {
	err := Map(DefGoroutines, data, function)
	if err != nil {
		panic(err)
	}
}

// Mapped calls function once for each item in sequence and returns a sequence with each mapped item as a result
func Mapped[S any, R any](goroutines int, data []S, function func(S) R) ([]R, error) {
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

// DefMapped calls function once for each item in sequence and returns a sequence with each mapped item as a result
func DefMapped[S any, R any](data []S, function func(S) R) ([]R, error) {
	result, err := Mapped(DefGoroutines, data, function)
	if err != nil {
		panic(err)
	}
	return result, nil
}

// MappedReduced calls mapFunction once for each item in sequence. The return value of each mapFunction is passed to reduceFunction.
//
// Note that while mapFunction is called concurrently, only one thread at a time will call reduceFunction.
// The order in which reduceFunction is called is determined by reduceOptions.
func MappedReduced[S any, R any](goroutines int, data []S, mapFunction func(S) R, reduceFunction func(*R, R), reduceOptions ReduceOption) (R, error) {
	var empty R
	if goroutines == 0 {
		return empty, fmt.Errorf("goroutines = 0")
	}
	if len(data) == 0 {
		return empty, fmt.Errorf("empty data")
	}

	switch reduceOptions {
	case UnorderedReduce:
		return unorderedMappedReduced(goroutines, data, mapFunction, reduceFunction)

	case OrderedReduce:
		return orderedMappedReduced(goroutines, data, mapFunction, reduceFunction)

	default:
		return empty, fmt.Errorf("invalid reduce order")
	}
}

func orderedMappedReduced[S any, R any](goroutines int, data []S, mapFunction func(S) R, reduceFunction func(*R, R)) (R, error) {
	result := make([]chan R, goroutines)
	for i := range result {
		result[i] = make(chan R, goroutines)
	}
	for i := 0; i < goroutines; i++ {
		go func(start, step int, datSize int) {
			for index := start; index < datSize; index += step {
				idx := index % goroutines
				result[idx] <- mapFunction(data[index])
			}
		}(i, goroutines, len(data))
	}

	accumulator := <-result[0]
	for j := 1; j < len(data); j++ {
		idx := j % goroutines
		res := <-result[idx]
		reduceFunction(&accumulator, res)
	}

	for _, ch := range result {
		close(ch)
	}

	return accumulator, nil
}

func unorderedMappedReduced[S any, R any](goroutines int, data []S, mapFunction func(S) R, reduceFunction func(*R, R)) (R, error) {
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	result := make(chan R, goroutines)
	for i := 0; i < goroutines; i++ {
		go func(start, step int, datSize int) {
			defer wg.Done()

			for index := start; index < datSize; index += step {
				result <- mapFunction(data[index])
			}
		}(i, goroutines, len(data))
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	accumulator := <-result
	for res := range result {
		reduceFunction(&accumulator, res)
	}

	return accumulator, nil
}
