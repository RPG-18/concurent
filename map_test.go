package concurent

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestInPlaceMap(t *testing.T) {
	t.Run("calculation", func(t *testing.T) {
		tests := []struct {
			data     []int
			expected []int
		}{{
			data:     []int{},
			expected: []int{},
		}, {
			data:     []int{1},
			expected: []int{4},
		}, {
			data:     []int{1, 1},
			expected: []int{4, 4},
		}, {
			data:     []int{1, 1, 1, 1, 1},
			expected: []int{4, 4, 4, 4, 4},
		}}

		for _, test := range tests {
			InPlaceMap(4, test.data, func(x *int) {
				*x = *x * 4
			})
			assert.Equal(t, test.expected, test.data)
		}
	})
	t.Run("concurrency", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5, 6, 7, 8}
		start := time.Now()
		InPlaceMap(4, data, func(x *int) {
			*x = *x * 2
			time.Sleep(time.Millisecond * 4)
		})
		lineTime := time.Millisecond * time.Duration(4*len(data))
		assert.True(t, time.Since(start) < lineTime)
	})
	t.Run("large", func(t *testing.T) {
		data := make([]int, 1000)
		InPlaceMap(3, data, func(x *int) {
			*x = 2
		})
		for _, v := range data {
			assert.Equal(t, 2, v)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("strconv", func(t *testing.T) {
		tests := []struct {
			data     []int
			expected []string
		}{{
			data:     []int{},
			expected: nil,
		}, {
			data:     []int{1},
			expected: []string{"1"},
		}, {
			data:     []int{1, 2},
			expected: []string{"1", "2"},
		}, {
			data:     []int{1, 2, 3, 4, 5, 6},
			expected: []string{"1", "2", "3", "4", "5", "6"},
		}}

		for _, test := range tests {
			result, err := Map(4, test.data, func(s int) string {
				return strconv.Itoa(s)
			})

			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	})
	t.Run("large", func(t *testing.T) {
		data := make([]int, 1000)
		for i := 0; i < len(data); i++ {
			data[i] = i
		}

		result, err := Map(4, data, func(s int) string {
			return strconv.Itoa(s)
		})
		assert.NoError(t, err)
		for i := 0; i < len(result); i++ {
			assert.Equal(t, strconv.Itoa(i), result[i])
		}
	})
}

func BenchmarkInPlaceMap(b *testing.B) {
	data := make([]int, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InPlaceMap(4, data, func(x *int) {
			*x = 100
			*x = *x * 2
		})
	}
}
