package convutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{nil, ""},
		{"hello", "hello"},
		{true, "true"},
		{false, "false"},
		{42, "42"},
		{int8(42), "42"},
		{int16(42), "42"},
		{int32(42), "42"},
		{int64(42), "42"},
		{uint(42), "42"},
		{uint8(42), "42"},
		{uint16(42), "42"},
		{uint32(42), "42"},
		{uint64(42), "42"},
		{float32(3.14), "3.14"},
		{float64(3.14), "3.14"},
	}

	for _, test := range tests {
		result := ToString(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    int
		expectError bool
	}{
		{nil, 0, false},
		{42, 42, false},
		{int8(42), 42, false},
		{int16(42), 42, false},
		{int32(42), 42, false},
		{int64(42), 42, false},
		{uint(42), 42, false},
		{uint8(42), 42, false},
		{uint16(42), 42, false},
		{uint32(42), 42, false},
		{uint64(42), 42, false},
		{float32(42.7), 42, false},
		{float64(42.7), 42, false},
		{true, 1, false},
		{false, 0, false},
		{"42", 42, false},
		{"invalid", 0, true},
		{[]int{1, 2, 3}, 0, true},
	}

	for _, test := range tests {
		result, err := ToInt(test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    int64
		expectError bool
	}{
		{nil, 0, false},
		{42, 42, false},
		{int64(42), 42, false},
		{float64(42.7), 42, false},
		{true, 1, false},
		{false, 0, false},
		{"42", 42, false},
		{"invalid", 0, true},
	}

	for _, test := range tests {
		result, err := ToInt64(test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    float64
		expectError bool
	}{
		{nil, 0, false},
		{42, 42.0, false},
		{float64(3.14), 3.14, false},
		{true, 1.0, false},
		{false, 0.0, false},
		{"3.14", 3.14, false},
		{"invalid", 0, true},
	}

	for _, test := range tests {
		result, err := ToFloat64(test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    bool
		expectError bool
	}{
		{nil, false, false},
		{true, true, false},
		{false, false, false},
		{1, true, false},
		{0, false, false},
		{42, true, false},
		{uint(1), true, false},
		{uint(0), false, false},
		{float64(1.0), true, false},
		{float64(0.0), false, false},
		{"true", true, false},
		{"false", false, false},
		{"1", true, false},
		{"0", false, false},
		{"invalid", false, true},
	}

	for _, test := range tests {
		result, err := ToBool(test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestToTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		input       interface{}
		expectError bool
	}{
		{nil, false},
		{now, false},
		{"2023-01-01T12:00:00Z", false},
		{"2023-01-01 12:00:00", false},
		{"2023-01-01", false},
		{now.Unix(), false},
		{"invalid", true},
		{42.5, true},
	}

	for _, test := range tests {
		result, err := ToTime(test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			if test.input == nil {
				assert.True(t, result.IsZero())
			} else {
				assert.False(t, result.IsZero())
			}
		}
	}
}

func TestToStringSlice(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    []string
		expectError bool
	}{
		{nil, nil, false},
		{[]int{1, 2, 3}, []string{"1", "2", "3"}, false},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, false},
		{[3]int{1, 2, 3}, []string{"1", "2", "3"}, false},
		{"not a slice", nil, true},
	}

	for _, test := range tests {
		result, err := ToStringSlice(test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestToIntSlice(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    []int
		expectError bool
	}{
		{nil, nil, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, false},
		{[]string{"1", "2", "3"}, []int{1, 2, 3}, false},
		{[]float64{1.7, 2.3, 3.9}, []int{1, 2, 3}, false},
		{[]string{"1", "invalid", "3"}, nil, true},
		{"not a slice", nil, true},
	}

	for _, test := range tests {
		result, err := ToIntSlice(test.input)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{map[string]interface{}{"name": "John", "age": 30}, `{"age":30,"name":"John"}`},
		{[]int{1, 2, 3}, `[1,2,3]`},
		{"hello", `"hello"`},
		{42, "42"},
		{true, "true"},
	}

	for _, test := range tests {
		result, err := ToJSON(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}

func TestFromJSON(t *testing.T) {
	var result map[string]interface{}
	err := FromJSON(`{"name": "John", "age": 30}`, &result)
	assert.NoError(t, err)
	assert.Equal(t, "John", result["name"])
	assert.Equal(t, float64(30), result["age"]) // JSON numbers are float64
}

func TestToMap(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string
	}

	person := Person{Name: "John", Age: 30, City: "New York"}
	result, err := ToMap(person)

	assert.NoError(t, err)
	assert.Equal(t, "John", result["name"])
	assert.Equal(t, 30, result["age"])
	assert.Equal(t, "New York", result["City"])
}

func TestToMap_WithPointer(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person := &Person{Name: "John", Age: 30}
	result, err := ToMap(person)

	assert.NoError(t, err)
	assert.Equal(t, "John", result["name"])
	assert.Equal(t, 30, result["age"])
}

func TestToMap_NonStruct(t *testing.T) {
	_, err := ToMap("not a struct")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a struct")
}

func TestMustFunctions(t *testing.T) {
	// Test successful conversions
	assert.Equal(t, "42", MustToString(42))
	assert.Equal(t, 42, MustToInt("42"))
	assert.Equal(t, int64(42), MustToInt64("42"))
	assert.Equal(t, 42.0, MustToFloat64("42"))
	assert.Equal(t, true, MustToBool("true"))

	// Test panics
	assert.Panics(t, func() { MustToInt("invalid") })
	assert.Panics(t, func() { MustToFloat64("invalid") })
	assert.Panics(t, func() { MustToBool("invalid") })
}

func TestMustToTime(t *testing.T) {
	// Test successful conversion
	result := MustToTime("2023-01-01T12:00:00Z")
	assert.False(t, result.IsZero())

	// Test panic
	assert.Panics(t, func() { MustToTime("invalid") })
}

func BenchmarkToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToString(42)
	}
}

func BenchmarkToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToInt("42")
	}
}

func BenchmarkToJSON(b *testing.B) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "New York",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToJSON(data)
	}
}
