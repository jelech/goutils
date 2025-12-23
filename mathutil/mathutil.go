// Package mathutil provides mathematical utilities for int and int64 types.
package mathutil

import (
	"math"
	"sort"
)

// MaxInt returns the maximum of two int values
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt returns the minimum of two int values
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt64 returns the maximum of two int64 values
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// MinInt64 returns the minimum of two int64 values
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MaxIntSlice returns the maximum value in an int slice
func MaxIntSlice(slice []int) int {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// MinIntSlice returns the minimum value in an int slice
func MinIntSlice(slice []int) int {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	min := slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// MaxInt64Slice returns the maximum value in an int64 slice
func MaxInt64Slice(slice []int64) int64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// MaxFloat64Slice returns the maximum value in an float64 slice
func MaxFloat64Slice(slice []float64) float64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// MaxFloat64Slice returns the maximum value in an float64 slice
func MaxFloat64Slice(slice []float64) float64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// MinFloat64Slice returns the minimum value in an float64 slice
func MinFloat64Slice(slice []float64) float64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	min := slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// MinInt64Slice returns the minimum value in an int64 slice
func MinInt64Slice(slice []int64) int64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	min := slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// MedianInt returns the median value of an int slice
func MedianInt(slice []int) float64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]int, len(slice))
	copy(sorted, slice)
	sort.Ints(sorted)

	n := len(sorted)
	if n%2 == 0 {
		// Even number of elements
		return float64(sorted[n/2-1]+sorted[n/2]) / 2.0
	}
	// Odd number of elements
	return float64(sorted[n/2])
}

// MedianInt64 returns the median value of an int64 slice
func MedianInt64(slice []int64) float64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]int64, len(slice))
	copy(sorted, slice)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	n := len(sorted)
	if n%2 == 0 {
		// Even number of elements
		return float64(sorted[n/2-1]+sorted[n/2]) / 2.0
	}
	// Odd number of elements
	return float64(sorted[n/2])
}

// ModeInt returns the most frequently occurring value(s) in an int slice
func ModeInt(slice []int) []int {
	if len(slice) == 0 {
		return nil
	}

	frequency := make(map[int]int)
	for _, v := range slice {
		frequency[v]++
	}

	maxFreq := 0
	for _, freq := range frequency {
		if freq > maxFreq {
			maxFreq = freq
		}
	}

	var modes []int
	for value, freq := range frequency {
		if freq == maxFreq {
			modes = append(modes, value)
		}
	}

	sort.Ints(modes)
	return modes
}

// ModeInt64 returns the most frequently occurring value(s) in an int64 slice
func ModeInt64(slice []int64) []int64 {
	if len(slice) == 0 {
		return nil
	}

	frequency := make(map[int64]int)
	for _, v := range slice {
		frequency[v]++
	}

	maxFreq := 0
	for _, freq := range frequency {
		if freq > maxFreq {
			maxFreq = freq
		}
	}

	var modes []int64
	for value, freq := range frequency {
		if freq == maxFreq {
			modes = append(modes, value)
		}
	}

	sort.Slice(modes, func(i, j int) bool {
		return modes[i] < modes[j]
	})
	return modes
}

// SumInt returns the sum of all values in an int slice
func SumInt(slice []int) int {
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return sum
}

// SumInt64 returns the sum of all values in an int64 slice
func SumInt64(slice []int64) int64 {
	var sum int64 = 0
	for _, v := range slice {
		sum += v
	}
	return sum
}

// AverageInt returns the average of all values in an int slice
func AverageInt(slice []int) float64 {
	if len(slice) == 0 {
		return 0
	}
	return float64(SumInt(slice)) / float64(len(slice))
}

// AverageInt64 returns the average of all values in an int64 slice
func AverageInt64(slice []int64) float64 {
	if len(slice) == 0 {
		return 0
	}
	return float64(SumInt64(slice)) / float64(len(slice))
}

// VarianceInt calculates the variance of an int slice
func VarianceInt(slice []int) float64 {
	if len(slice) == 0 {
		return 0
	}

	mean := AverageInt(slice)
	var variance float64
	for _, v := range slice {
		diff := float64(v) - mean
		variance += diff * diff
	}
	return variance / float64(len(slice))
}

// VarianceInt64 calculates the variance of an int64 slice
func VarianceInt64(slice []int64) float64 {
	if len(slice) == 0 {
		return 0
	}

	mean := AverageInt64(slice)
	var variance float64
	for _, v := range slice {
		diff := float64(v) - mean
		variance += diff * diff
	}
	return variance / float64(len(slice))
}

// StandardDeviationInt calculates the standard deviation of an int slice
func StandardDeviationInt(slice []int) float64 {
	return math.Sqrt(VarianceInt(slice))
}

// StandardDeviationInt64 calculates the standard deviation of an int64 slice
func StandardDeviationInt64(slice []int64) float64 {
	return math.Sqrt(VarianceInt64(slice))
}

// PercentileInt returns the value at the given percentile (0-100) in an int slice
func PercentileInt(slice []int, percentile float64) float64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	if percentile < 0 || percentile > 100 {
		panic("percentile must be between 0 and 100")
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]int, len(slice))
	copy(sorted, slice)
	sort.Ints(sorted)

	if percentile == 100 {
		return float64(sorted[len(sorted)-1])
	}

	index := percentile / 100.0 * float64(len(sorted)-1)
	lower := int(index)
	upper := lower + 1

	if upper >= len(sorted) {
		return float64(sorted[lower])
	}

	weight := index - float64(lower)
	return float64(sorted[lower])*(1-weight) + float64(sorted[upper])*weight
}

// PercentileInt64 returns the value at the given percentile (0-100) in an int64 slice
func PercentileInt64(slice []int64, percentile float64) float64 {
	if len(slice) == 0 {
		panic("slice cannot be empty")
	}
	if percentile < 0 || percentile > 100 {
		panic("percentile must be between 0 and 100")
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]int64, len(slice))
	copy(sorted, slice)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	if percentile == 100 {
		return float64(sorted[len(sorted)-1])
	}

	index := percentile / 100.0 * float64(len(sorted)-1)
	lower := int(index)
	upper := lower + 1

	if upper >= len(sorted) {
		return float64(sorted[lower])
	}

	weight := index - float64(lower)
	return float64(sorted[lower])*(1-weight) + float64(sorted[upper])*weight
}

// RangeInt returns the difference between max and min values in an int slice
func RangeInt(slice []int) int {
	if len(slice) == 0 {
		return 0
	}
	return MaxIntSlice(slice) - MinIntSlice(slice)
}

// RangeInt64 returns the difference between max and min values in an int64 slice
func RangeInt64(slice []int64) int64 {
	if len(slice) == 0 {
		return 0
	}
	return MaxInt64Slice(slice) - MinInt64Slice(slice)
}

// UniqueInt returns unique values from an int slice
func UniqueInt(slice []int) []int {
	seen := make(map[int]bool)
	var unique []int

	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}

// UniqueInt64 returns unique values from an int64 slice
func UniqueInt64(slice []int64) []int64 {
	seen := make(map[int64]bool)
	var unique []int64

	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}

// CountInt counts occurrences of a specific value in an int slice
func CountInt(slice []int, value int) int {
	count := 0
	for _, v := range slice {
		if v == value {
			count++
		}
	}
	return count
}

// CountInt64 counts occurrences of a specific value in an int64 slice
func CountInt64(slice []int64, value int64) int {
	count := 0
	for _, v := range slice {
		if v == value {
			count++
		}
	}
	return count
}

// ClampInt clamps an int value between min and max
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt64 clamps an int64 value between min and max
func ClampInt64(value, min, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// GCD returns the greatest common divisor of two integers
func GCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// GCDInt64 returns the greatest common divisor of two int64 values
func GCDInt64(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple of two integers
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	return (a * b) / GCD(a, b)
}

// LCMInt64 returns the least common multiple of two int64 values
func LCMInt64(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	return (a * b) / GCDInt64(a, b)
}

// IsPowerOfTwo checks if a number is a power of two
func IsPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}

// IsPowerOfTwoInt64 checks if an int64 number is a power of two
func IsPowerOfTwoInt64(n int64) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}
