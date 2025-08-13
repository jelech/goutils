# mathutil

mathutil 包提供了针对 int 和 int64 类型的数学工具函数，包括基本比较、统计计算、数组操作等功能。

## 功能特性

### 基本比较函数
- `MaxInt(a, b int) int` - 返回两个int值的最大值
- `MinInt(a, b int) int` - 返回两个int值的最小值  
- `MaxInt64(a, b int64) int64` - 返回两个int64值的最大值
- `MinInt64(a, b int64) int64` - 返回两个int64值的最小值

### 数组操作函数
- `MaxIntSlice(slice []int) int` - 返回int切片中的最大值
- `MinIntSlice(slice []int) int` - 返回int切片中的最小值
- `MaxInt64Slice(slice []int64) int64` - 返回int64切片中的最大值
- `MinInt64Slice(slice []int64) int64` - 返回int64切片中的最小值

### 统计计算函数
- `SumInt(slice []int) int` - 计算int切片的总和
- `SumInt64(slice []int64) int64` - 计算int64切片的总和
- `AverageInt(slice []int) float64` - 计算int切片的平均值
- `AverageInt64(slice []int64) float64` - 计算int64切片的平均值
- `MedianInt(slice []int) float64` - 计算int切片的中位数
- `MedianInt64(slice []int64) float64` - 计算int64切片的中位数
- `ModeInt(slice []int) []int` - 返回int切片中的众数
- `ModeInt64(slice []int64) []int64` - 返回int64切片中的众数
- `VarianceInt(slice []int) float64` - 计算int切片的方差
- `VarianceInt64(slice []int64) float64` - 计算int64切片的方差
- `StandardDeviationInt(slice []int) float64` - 计算int切片的标准差
- `StandardDeviationInt64(slice []int64) float64` - 计算int64切片的标准差
- `PercentileInt(slice []int, percentile float64) float64` - 计算int切片的百分位数
- `PercentileInt64(slice []int64, percentile float64) float64` - 计算int64切片的百分位数
- `RangeInt(slice []int) int` - 计算int切片的极差（最大值-最小值）
- `RangeInt64(slice []int64) int64` - 计算int64切片的极差

### 数组工具函数
- `UniqueInt(slice []int) []int` - 返回int切片中的唯一值
- `UniqueInt64(slice []int64) []int64` - 返回int64切片中的唯一值
- `CountInt(slice []int, value int) int` - 计算int切片中特定值的出现次数
- `CountInt64(slice []int64, value int64) int` - 计算int64切片中特定值的出现次数

### 数学工具函数
- `ClampInt(value, min, max int) int` - 将int值限制在指定范围内
- `ClampInt64(value, min, max int64) int64` - 将int64值限制在指定范围内
- `GCD(a, b int) int` - 计算两个int值的最大公约数
- `GCDInt64(a, b int64) int64` - 计算两个int64值的最大公约数
- `LCM(a, b int) int` - 计算两个int值的最小公倍数
- `LCMInt64(a, b int64) int64` - 计算两个int64值的最小公倍数
- `IsPowerOfTwo(n int) bool` - 检查int值是否为2的幂
- `IsPowerOfTwoInt64(n int64) bool` - 检查int64值是否为2的幂

## 使用示例

```go
package main

import (
    "fmt"
    "github.com/jelech/goutils/mathutil"
)

func main() {
    // 基本比较
    max := mathutil.MaxInt(5, 3) // 返回 5
    min := mathutil.MinInt64(10, 7) // 返回 7
    
    // 数组操作
    numbers := []int{4, 2, 7, 1, 9, 3}
    maxVal := mathutil.MaxIntSlice(numbers) // 返回 9
    sum := mathutil.SumInt(numbers) // 返回 26
    avg := mathutil.AverageInt(numbers) // 返回 4.33
    median := mathutil.MedianInt(numbers) // 返回 3.5
    
    // 统计计算
    data := []int{1, 2, 2, 3, 4, 4, 4, 5}
    mode := mathutil.ModeInt(data) // 返回 [4]
    variance := mathutil.VarianceInt(data) // 计算方差
    stdDev := mathutil.StandardDeviationInt(data) // 计算标准差
    p75 := mathutil.PercentileInt(data, 75) // 计算75百分位数
    
    // 数学工具
    gcd := mathutil.GCD(12, 18) // 返回 6
    lcm := mathutil.LCM(12, 18) // 返回 36
    isPower := mathutil.IsPowerOfTwo(16) // 返回 true
    clamped := mathutil.ClampInt(15, 1, 10) // 返回 10
}
```

## 设计原则

1. **兼容性**: 为了兼容Go 1.17，没有使用泛型，而是为int和int64类型分别提供函数
2. **安全性**: 对空切片等边界情况进行适当处理
3. **性能**: 避免不必要的内存分配，统计函数会创建数据副本以避免修改原始数据
4. **一致性**: 函数命名和行为保持一致性，遵循Go语言惯例

## 注意事项

- 大部分切片操作函数在遇到空切片时会panic，使用前请检查切片长度
- 中位数、百分位数等函数会创建数据的副本进行排序，不会修改原始数据
- 浮点数相关的计算可能存在精度问题，建议根据具体需求进行适当的舍入处理
