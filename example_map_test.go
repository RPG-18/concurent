package concurrent_test

import (
	"fmt"
	"strconv"

	concurrent "github.com/RPG-18/concurrent"
)

func ExampleMapped() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	strings, _ := concurrent.Mapped(4, numbers, func(value int) string {
		return strconv.Itoa(value)
	})

	fmt.Println(strings)
}

func ExampleMap() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	_ = concurrent.Map(4, numbers, func(value *int) {
		*value *= 2
	})

	fmt.Println(numbers)
}

func ExampleMappedReduced() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	ordered, _ := concurrent.MappedReduced(3, numbers, func(value int) string {
		return strconv.Itoa(value)
	}, func(sum *string, value string) {
		*sum = *sum + value
	}, concurrent.OrderedReduce)
	unordered, _ := concurrent.MappedReduced(3, numbers, func(value int) string {
		return strconv.Itoa(value)
	}, func(sum *string, value string) {
		*sum = *sum + value
	}, concurrent.UnorderedReduce)

	fmt.Println(ordered, unordered) // 123456 142536
}
