package concurrent

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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
			_ = Map(4, test.data, func(x *int) {
				*x = *x * 4
			})
			assert.Equal(t, test.expected, test.data)
		}
	})
	t.Run("large", func(t *testing.T) {
		data := make([]int, 1000)
		_ = Map(3, data, func(x *int) {
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
			result, err := Mapped(4, test.data, func(s int) string {
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

		result, err := Mapped(4, data, func(s int) string {
			return strconv.Itoa(s)
		})
		assert.NoError(t, err)
		for i := 0; i < len(result); i++ {
			assert.Equal(t, strconv.Itoa(i), result[i])
		}
	})
}

func TestDefInPlaceMap(t *testing.T) {
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
			DefMap(test.data, func(x *int) {
				*x = *x * 4
			})
			assert.Equal(t, test.expected, test.data)
		}
	})
}

func TestDefMap(t *testing.T) {
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
			result, err := DefMapped(test.data, func(s int) string {
				return strconv.Itoa(s)
			})

			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	})
}

func TestMappedReduced(t *testing.T) {
	t.Run("ordered", func(t *testing.T) {
		t.Run("result int", func(t *testing.T) {
			result, err := MappedReduced(3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(x int) int {
				return x * 2
			}, func(r *int, x int) {
				*r = *r + x
			}, OrderedReduce)

			ss := 0
			for _, v := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
				ss = ss + v*2
			}
			assert.NoError(t, err)
			assert.Equal(t, ss, result)
		})
		t.Run("result string", func(t *testing.T) {
			result, err := MappedReduced(3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(x int) string {
				return strconv.Itoa(x)
			}, func(r *string, x string) {
				*r = *r + x
			}, OrderedReduce)
			assert.NoError(t, err)
			assert.Equal(t, "123456789", result)
		})
	})

	t.Run("unordered", func(t *testing.T) {
		t.Run("result int", func(t *testing.T) {
			result, err := MappedReduced(3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(x int) int {
				return x * 2
			}, func(r *int, x int) {
				*r = *r + x
			}, UnorderedReduce)

			ss := 0
			for _, v := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
				ss = ss + v*2
			}
			assert.NoError(t, err)
			assert.Equal(t, ss, result)
		})
		t.Run("result string", func(t *testing.T) {
			result, err := MappedReduced(3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(x int) string {
				return strconv.Itoa(x)
			}, func(r *string, x string) {
				*r = *r + x
			}, UnorderedReduce)
			assert.NoError(t, err)
			assert.Len(t, result, 9)
			assert.NotEqual(t, "123456789", result)
		})
	})
}

func BenchmarkInPlaceMap(b *testing.B) {
	data := make([]int, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(4, data, func(x *int) {
			*x = 100
			*x = *x * 2
		})
	}
}

func BenchmarkMappedReduced(b *testing.B) {
	data := make([]int, 1024*12)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()

	b.Run("ordered", func(b *testing.B) {
		MappedReduced(6, data, func(x int) int {
			return x * 2
		}, func(r *int, x int) {
			*r = +x
		}, OrderedReduce)
	})
	b.Run("unordered", func(b *testing.B) {
		MappedReduced(6, data, func(x int) int {
			return x * 2
		}, func(r *int, x int) {
			*r = +x
		}, UnorderedReduce)
	})
}
