# Simple MapReduce on Go

[![Test](https://github.com/RPG-18/concurent/actions/workflows/concurrent.yml/badge.svg?branch=main)](https://github.com/RPG-18/concurrent/actions)
[![GoDoc](https://pkg.go.dev/badge/github.com/RPG-18/concurrent?status.svg)](https://pkg.go.dev/github.com/RPG-18/concurrent?tab=doc)


Package concurrent provides high-level APIs for concurrent programming. Package implemented [Qt Concurrent API](https://doc.qt.io/qt-6/qtconcurrent-index.html)

## Examples

#### Map
```go
package main

import (
	"fmt"
	concurrent "github.com/RPG-18/concurrent"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	_ = concurrent.Map(4, numbers, func(value *int) {
		*value *= 2
	})

	fmt.Println(numbers) // [2 4 6 8 10 12]
}
```

#### Mapped
```go
package main

import (
	"fmt"
	"strconv"

	concurrent "github.com/RPG-18/concurrent"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	strings, _ := concurrent.Mapped(4, numbers, func(value int) string {
		return strconv.Itoa(value)
	})

	fmt.Println(strings) // [1 2 3 4 5 6]
}
```

#### MappedReduced
```go
package main

import (
	"fmt"
	"strconv"

	concurrent "github.com/RPG-18/concurrent"
)

func main() {
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
```