package concurrent_test

import (
	"fmt"
	"strconv"

	concurrent "github.com/RPG-18/concurrent"
)

func ExampleMap() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	strings, _ := concurrent.Map(4, numbers, func(value int) string {
		return strconv.Itoa(value)
	})

	fmt.Println(strings)
}

func ExampleInPlaceMap() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	_ = concurrent.InPlaceMap(4, numbers, func(value *int) {
		*value *= 2
	})

	fmt.Println(numbers)
}
