package mathutil

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 5, MaxInt(3, 5))
	assert.Equal(t, 5, MaxInt(5, 3))
	assert.Equal(t, 0, MaxInt(-1, 0))
	assert.Equal(t, -1, MaxInt(-5, -1))
}

func TestMinInt(t *testing.T) {
	assert.Equal(t, 3, MinInt(3, 5))
	assert.Equal(t, 3, MinInt(5, 3))
	assert.Equal(t, -1, MinInt(-1, 0))
	assert.Equal(t, -5, MinInt(-5, -1))
}

func TestMaxInt64(t *testing.T) {
	assert.Equal(t, int64(5), MaxInt64(3, 5))
	assert.Equal(t, int64(5), MaxInt64(5, 3))
	assert.Equal(t, int64(0), MaxInt64(-1, 0))
	assert.Equal(t, int64(-1), MaxInt64(-5, -1))
}

func TestMinInt64(t *testing.T) {
	assert.Equal(t, int64(3), MinInt64(3, 5))
	assert.Equal(t, int64(3), MinInt64(5, 3))
	assert.Equal(t, int64(-1), MinInt64(-1, 0))
	assert.Equal(t, int64(-5), MinInt64(-5, -1))
}

func TestMaxIntSlice(t *testing.T) {
	assert.Equal(t, 5, MaxIntSlice([]int{1, 3, 5, 2}))
	assert.Equal(t, 10, MaxIntSlice([]int{10}))
	assert.Equal(t, -1, MaxIntSlice([]int{-5, -3, -1}))

	// Test panic for empty slice
	assert.Panics(t, func() {
		MaxIntSlice([]int{})
	})
}

func TestMinIntSlice(t *testing.T) {
	assert.Equal(t, 1, MinIntSlice([]int{1, 3, 5, 2}))
	assert.Equal(t, 10, MinIntSlice([]int{10}))
	assert.Equal(t, -5, MinIntSlice([]int{-5, -3, -1}))

	// Test panic for empty slice
	assert.Panics(t, func() {
		MinIntSlice([]int{})
	})
}

func TestMaxInt64Slice(t *testing.T) {
	assert.Equal(t, int64(5), MaxInt64Slice([]int64{1, 3, 5, 2}))
	assert.Equal(t, int64(10), MaxInt64Slice([]int64{10}))
	assert.Equal(t, int64(-1), MaxInt64Slice([]int64{-5, -3, -1}))

	// Test panic for empty slice
	assert.Panics(t, func() {
		MaxInt64Slice([]int64{})
	})
}

func TestMinInt64Slice(t *testing.T) {
	assert.Equal(t, int64(1), MinInt64Slice([]int64{1, 3, 5, 2}))
	assert.Equal(t, int64(10), MinInt64Slice([]int64{10}))
	assert.Equal(t, int64(-5), MinInt64Slice([]int64{-5, -3, -1}))

	// Test panic for empty slice
	assert.Panics(t, func() {
		MinInt64Slice([]int64{})
	})
}

func TestMedianInt(t *testing.T) {
	// Odd number of elements
	assert.Equal(t, 3.0, MedianInt([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 3.0, MedianInt([]int{5, 1, 3, 2, 4}))

	// Even number of elements
	assert.Equal(t, 2.5, MedianInt([]int{1, 2, 3, 4}))
	assert.Equal(t, 2.5, MedianInt([]int{4, 1, 3, 2}))

	// Single element
	assert.Equal(t, 5.0, MedianInt([]int{5}))

	// Test panic for empty slice
	assert.Panics(t, func() {
		MedianInt([]int{})
	})
}

func TestMedianInt64(t *testing.T) {
	// Odd number of elements
	assert.Equal(t, 3.0, MedianInt64([]int64{1, 2, 3, 4, 5}))
	assert.Equal(t, 3.0, MedianInt64([]int64{5, 1, 3, 2, 4}))

	// Even number of elements
	assert.Equal(t, 2.5, MedianInt64([]int64{1, 2, 3, 4}))
	assert.Equal(t, 2.5, MedianInt64([]int64{4, 1, 3, 2}))

	// Single element
	assert.Equal(t, 5.0, MedianInt64([]int64{5}))

	// Test panic for empty slice
	assert.Panics(t, func() {
		MedianInt64([]int64{})
	})
}

func TestModeInt(t *testing.T) {
	// Single mode
	assert.Equal(t, []int{2}, ModeInt([]int{1, 2, 2, 3}))

	// Multiple modes
	modes := ModeInt([]int{1, 1, 2, 2, 3})
	assert.Contains(t, modes, 1)
	assert.Contains(t, modes, 2)
	assert.Len(t, modes, 2)

	// All elements same frequency
	modes = ModeInt([]int{1, 2, 3})
	assert.Len(t, modes, 3)

	// Empty slice
	assert.Nil(t, ModeInt([]int{}))
}

func TestModeInt64(t *testing.T) {
	// Single mode
	assert.Equal(t, []int64{2}, ModeInt64([]int64{1, 2, 2, 3}))

	// Multiple modes
	modes := ModeInt64([]int64{1, 1, 2, 2, 3})
	assert.Contains(t, modes, int64(1))
	assert.Contains(t, modes, int64(2))
	assert.Len(t, modes, 2)

	// Empty slice
	assert.Nil(t, ModeInt64([]int64{}))
}

func TestSumInt(t *testing.T) {
	assert.Equal(t, 15, SumInt([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 10, SumInt([]int{10}))
	assert.Equal(t, 0, SumInt([]int{}))
	assert.Equal(t, -5, SumInt([]int{-1, -2, -2}))
}

func TestSumInt64(t *testing.T) {
	assert.Equal(t, int64(15), SumInt64([]int64{1, 2, 3, 4, 5}))
	assert.Equal(t, int64(10), SumInt64([]int64{10}))
	assert.Equal(t, int64(0), SumInt64([]int64{}))
	assert.Equal(t, int64(-5), SumInt64([]int64{-1, -2, -2}))
}

func TestAverageInt(t *testing.T) {
	assert.Equal(t, 3.0, AverageInt([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 10.0, AverageInt([]int{10}))
	assert.Equal(t, 0.0, AverageInt([]int{}))
	assert.InDelta(t, -1.667, AverageInt([]int{-1, -2, -2}), 0.001)
}

func TestAverageInt64(t *testing.T) {
	assert.Equal(t, 3.0, AverageInt64([]int64{1, 2, 3, 4, 5}))
	assert.Equal(t, 10.0, AverageInt64([]int64{10}))
	assert.Equal(t, 0.0, AverageInt64([]int64{}))
	assert.InDelta(t, -1.667, AverageInt64([]int64{-1, -2, -2}), 0.001)
}

func TestVarianceInt(t *testing.T) {
	// Known variance: [1,2,3,4,5] has variance = 2.0
	assert.Equal(t, 2.0, VarianceInt([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 0.0, VarianceInt([]int{5}))
	assert.Equal(t, 0.0, VarianceInt([]int{}))
}

func TestVarianceInt64(t *testing.T) {
	// Known variance: [1,2,3,4,5] has variance = 2.0
	assert.Equal(t, 2.0, VarianceInt64([]int64{1, 2, 3, 4, 5}))
	assert.Equal(t, 0.0, VarianceInt64([]int64{5}))
	assert.Equal(t, 0.0, VarianceInt64([]int64{}))
}

func TestStandardDeviationInt(t *testing.T) {
	// Standard deviation is sqrt of variance
	assert.InDelta(t, math.Sqrt(2.0), StandardDeviationInt([]int{1, 2, 3, 4, 5}), 0.001)
	assert.Equal(t, 0.0, StandardDeviationInt([]int{5}))
}

func TestStandardDeviationInt64(t *testing.T) {
	// Standard deviation is sqrt of variance
	assert.InDelta(t, math.Sqrt(2.0), StandardDeviationInt64([]int64{1, 2, 3, 4, 5}), 0.001)
	assert.Equal(t, 0.0, StandardDeviationInt64([]int64{5}))
}

func TestPercentileInt(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	assert.Equal(t, 1.0, PercentileInt(data, 0))
	assert.Equal(t, 5.5, PercentileInt(data, 50))
	assert.Equal(t, 10.0, PercentileInt(data, 100))

	// Test panic for empty slice
	assert.Panics(t, func() {
		PercentileInt([]int{}, 50)
	})

	// Test panic for invalid percentile
	assert.Panics(t, func() {
		PercentileInt(data, -1)
	})
	assert.Panics(t, func() {
		PercentileInt(data, 101)
	})
}

func TestPercentileInt64(t *testing.T) {
	data := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	assert.Equal(t, 1.0, PercentileInt64(data, 0))
	assert.Equal(t, 5.5, PercentileInt64(data, 50))
	assert.Equal(t, 10.0, PercentileInt64(data, 100))

	// Test panic for empty slice
	assert.Panics(t, func() {
		PercentileInt64([]int64{}, 50)
	})
}

func TestRangeInt(t *testing.T) {
	assert.Equal(t, 4, RangeInt([]int{1, 3, 5, 2}))
	assert.Equal(t, 0, RangeInt([]int{5}))
	assert.Equal(t, 0, RangeInt([]int{}))
	assert.Equal(t, 4, RangeInt([]int{-5, -3, -1}))
}

func TestRangeInt64(t *testing.T) {
	assert.Equal(t, int64(4), RangeInt64([]int64{1, 3, 5, 2}))
	assert.Equal(t, int64(0), RangeInt64([]int64{5}))
	assert.Equal(t, int64(0), RangeInt64([]int64{}))
	assert.Equal(t, int64(4), RangeInt64([]int64{-5, -3, -1}))
}

func TestUniqueInt(t *testing.T) {
	unique := UniqueInt([]int{1, 2, 2, 3, 3, 3})
	assert.Len(t, unique, 3)
	assert.Contains(t, unique, 1)
	assert.Contains(t, unique, 2)
	assert.Contains(t, unique, 3)

	assert.Nil(t, UniqueInt([]int{}))
	assert.Equal(t, []int{1}, UniqueInt([]int{1}))
}

func TestUniqueInt64(t *testing.T) {
	unique := UniqueInt64([]int64{1, 2, 2, 3, 3, 3})
	assert.Len(t, unique, 3)
	assert.Contains(t, unique, int64(1))
	assert.Contains(t, unique, int64(2))
	assert.Contains(t, unique, int64(3))

	assert.Nil(t, UniqueInt64([]int64{}))
	assert.Equal(t, []int64{1}, UniqueInt64([]int64{1}))
}

func TestCountInt(t *testing.T) {
	assert.Equal(t, 3, CountInt([]int{1, 2, 2, 2, 3}, 2))
	assert.Equal(t, 1, CountInt([]int{1, 2, 3}, 1))
	assert.Equal(t, 0, CountInt([]int{1, 2, 3}, 4))
	assert.Equal(t, 0, CountInt([]int{}, 1))
}

func TestCountInt64(t *testing.T) {
	assert.Equal(t, 3, CountInt64([]int64{1, 2, 2, 2, 3}, 2))
	assert.Equal(t, 1, CountInt64([]int64{1, 2, 3}, 1))
	assert.Equal(t, 0, CountInt64([]int64{1, 2, 3}, 4))
	assert.Equal(t, 0, CountInt64([]int64{}, 1))
}

func TestClampInt(t *testing.T) {
	assert.Equal(t, 5, ClampInt(10, 1, 5))
	assert.Equal(t, 1, ClampInt(-5, 1, 5))
	assert.Equal(t, 3, ClampInt(3, 1, 5))
	assert.Equal(t, 1, ClampInt(1, 1, 5))
	assert.Equal(t, 5, ClampInt(5, 1, 5))
}

func TestClampInt64(t *testing.T) {
	assert.Equal(t, int64(5), ClampInt64(10, 1, 5))
	assert.Equal(t, int64(1), ClampInt64(-5, 1, 5))
	assert.Equal(t, int64(3), ClampInt64(3, 1, 5))
	assert.Equal(t, int64(1), ClampInt64(1, 1, 5))
	assert.Equal(t, int64(5), ClampInt64(5, 1, 5))
}

func TestGCD(t *testing.T) {
	assert.Equal(t, 6, GCD(12, 18))
	assert.Equal(t, 1, GCD(17, 19))
	assert.Equal(t, 5, GCD(10, 15))
	assert.Equal(t, 7, GCD(14, 21))
	assert.Equal(t, 6, GCD(-12, 18))
	assert.Equal(t, 6, GCD(12, -18))
	assert.Equal(t, 6, GCD(-12, -18))
}

func TestGCDInt64(t *testing.T) {
	assert.Equal(t, int64(6), GCDInt64(12, 18))
	assert.Equal(t, int64(1), GCDInt64(17, 19))
	assert.Equal(t, int64(5), GCDInt64(10, 15))
	assert.Equal(t, int64(7), GCDInt64(14, 21))
	assert.Equal(t, int64(6), GCDInt64(-12, 18))
	assert.Equal(t, int64(6), GCDInt64(12, -18))
	assert.Equal(t, int64(6), GCDInt64(-12, -18))
}

func TestLCM(t *testing.T) {
	assert.Equal(t, 36, LCM(12, 18))
	assert.Equal(t, 323, LCM(17, 19))
	assert.Equal(t, 30, LCM(10, 15))
	assert.Equal(t, 42, LCM(14, 21))
	assert.Equal(t, 0, LCM(0, 5))
	assert.Equal(t, 0, LCM(5, 0))
}

func TestLCMInt64(t *testing.T) {
	assert.Equal(t, int64(36), LCMInt64(12, 18))
	assert.Equal(t, int64(323), LCMInt64(17, 19))
	assert.Equal(t, int64(30), LCMInt64(10, 15))
	assert.Equal(t, int64(42), LCMInt64(14, 21))
	assert.Equal(t, int64(0), LCMInt64(0, 5))
	assert.Equal(t, int64(0), LCMInt64(5, 0))
}

func TestIsPowerOfTwo(t *testing.T) {
	assert.True(t, IsPowerOfTwo(1))
	assert.True(t, IsPowerOfTwo(2))
	assert.True(t, IsPowerOfTwo(4))
	assert.True(t, IsPowerOfTwo(8))
	assert.True(t, IsPowerOfTwo(16))
	assert.True(t, IsPowerOfTwo(32))
	assert.True(t, IsPowerOfTwo(64))

	assert.False(t, IsPowerOfTwo(0))
	assert.False(t, IsPowerOfTwo(-1))
	assert.False(t, IsPowerOfTwo(3))
	assert.False(t, IsPowerOfTwo(5))
	assert.False(t, IsPowerOfTwo(6))
	assert.False(t, IsPowerOfTwo(7))
	assert.False(t, IsPowerOfTwo(9))
	assert.False(t, IsPowerOfTwo(15))
}

func TestIsPowerOfTwoInt64(t *testing.T) {
	assert.True(t, IsPowerOfTwoInt64(1))
	assert.True(t, IsPowerOfTwoInt64(2))
	assert.True(t, IsPowerOfTwoInt64(4))
	assert.True(t, IsPowerOfTwoInt64(8))
	assert.True(t, IsPowerOfTwoInt64(16))
	assert.True(t, IsPowerOfTwoInt64(32))
	assert.True(t, IsPowerOfTwoInt64(64))

	assert.False(t, IsPowerOfTwoInt64(0))
	assert.False(t, IsPowerOfTwoInt64(-1))
	assert.False(t, IsPowerOfTwoInt64(3))
	assert.False(t, IsPowerOfTwoInt64(5))
	assert.False(t, IsPowerOfTwoInt64(6))
	assert.False(t, IsPowerOfTwoInt64(7))
	assert.False(t, IsPowerOfTwoInt64(9))
	assert.False(t, IsPowerOfTwoInt64(15))
}

func TestNoiseGenerator(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== 中心值概率分布测试 ===")

	// 测试以4为中心的分布
	centerValue := 4
	length := 3
	probs := generateCenteredProbabilityDistribution(centerValue, length, 0.9)

	fmt.Printf("以%d为中心，控制因子0.8的概率分布:\n", centerValue)
	for i := 1; i <= 10; i++ {
		if probs[i] > 0.001 {
			fmt.Printf("P(%d) = %.1f%%\n", i, probs[i]*100)
		}
	}

	fmt.Println("\n=== 概率偏移工具测试 ===")

	// 测试概率偏移工具
	inputValue := 5
	controlFactors := []float64{0.1, 0.6, 0.9}

	for _, factor := range controlFactors {
		fmt.Printf("\n输入值: %d, 控制因子: %.1f\n", inputValue, factor)
		fmt.Printf("偏移结果: ")

		// 多次采样看分布效果
		results := make(map[int]int)
		for i := 0; i < 1000; i++ {
			shifted := applyProbabilityShift(inputValue, 3, factor)
			results[shifted]++
		}

		// 显示采样统计
		for value := 1; value <= 10; value++ {
			if count := results[value]; count > 0 {
				fmt.Printf("%d:%.0f%% ", value, float64(count)/10)
			}
		}
		fmt.Println()
	}

	fmt.Println("\n=== 单次偏移示例 ===")

	// 单次偏移示例
	testValues := []int{2, 5, 8}
	for _, val := range testValues {
		shifted := applyProbabilityShift(val, 3, 0.7)
		fmt.Printf("输入: %d → 偏移后: %d\n", val, shifted)
	}
}
