package mathutil

import (
	"math"
	"math/rand"
)

// 生成以指定值为中心的概率分布
func generateCenteredProbabilityDistribution(centerValue, length int, controlFactor float64) map[int]float64 {
	if controlFactor <= 0 || controlFactor > 1 {
		panic("control_factor must be in range (0, 1]")
	}

	minValue := MaxInt(centerValue-length, 0)
	maxValue := MinInt(centerValue+length, 10000)

	probabilities := make(map[int]float64)

	// 当控制因子为1时，100%概率在中心值
	if controlFactor == 1.0 {
		for i := minValue; i <= maxValue; i++ {
			if i == centerValue {
				probabilities[i] = 1.0
			} else {
				probabilities[i] = 0.0
			}
		}
		return probabilities
	}

	decayRate := -math.Log(1 - controlFactor + 0.001)
	weights := make(map[int]float64)
	totalWeight := 0.0

	// 计算每个值相对于中心值的距离，并计算权重
	for i := minValue; i <= maxValue; i++ {
		distance := math.Abs(float64(i - centerValue))
		weight := math.Exp(-decayRate * distance)
		weights[i] = weight
		totalWeight += weight
	}

	// 归一化得到概率
	for i := minValue; i <= maxValue; i++ {
		probabilities[i] = weights[i] / totalWeight
	}

	return probabilities
}

// 从概率分布中采样单个值
func sampleSingle(probabilities map[int]float64) int {
	r := rand.Float64()
	cumulative := 0.0

	// 按键值排序遍历
	for value := getMinKey(probabilities); value <= getMaxKey(probabilities); value++ {
		if prob, exists := probabilities[value]; exists {
			cumulative += prob
			if r <= cumulative {
				return value
			}
		}
	}

	// 兜底返回最大键值
	return getMaxKey(probabilities)
}

// 概率偏移工具函数：对输入数字进行概率偏移
func applyProbabilityShift(inputValue, length int, controlFactor float64) int {
	// 生成以输入值为中心的概率分布
	probs := generateCenteredProbabilityDistribution(inputValue, length, controlFactor)

	// 采样得到偏移后的值
	return sampleSingle(probs)
}

// 辅助函数：获取map中的最小键
func getMinKey(m map[int]float64) int {
	min := math.MaxInt32
	for k := range m {
		if k < min {
			min = k
		}
	}
	return min
}

// 辅助函数：获取map中的最大键
func getMaxKey(m map[int]float64) int {
	max := math.MinInt32
	for k := range m {
		if k > max {
			max = k
		}
	}
	return max
}
