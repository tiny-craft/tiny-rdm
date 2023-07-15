package mathutil

import (
	"math"
	. "tinyrdm/backend/utils"
)

// MaxWithIndex 查找所有元素中的最大值
func MaxWithIndex[T Hashable](items ...T) (T, int) {
	selIndex := -1
	for i, t := range items {
		if selIndex < 0 {
			selIndex = i
		} else {
			if t > items[selIndex] {
				selIndex = i
			}
		}
	}
	return items[selIndex], selIndex
}

func Max[T Hashable](items ...T) T {
	val, _ := MaxWithIndex(items...)
	return val
}

// MinWithIndex 查找所有元素中的最小值
func MinWithIndex[T Hashable](items ...T) (T, int) {
	selIndex := -1
	for i, t := range items {
		if selIndex < 0 {
			selIndex = i
		} else {
			if t < items[selIndex] {
				selIndex = i
			}
		}
	}
	return items[selIndex], selIndex
}

func Min[T Hashable](items ...T) T {
	val, _ := MinWithIndex(items...)
	return val
}

// Clamp 返回限制在minVal和maxVal范围内的value
func Clamp[T Hashable](value T, minVal T, maxVal T) T {
	if minVal > maxVal {
		minVal, maxVal = maxVal, minVal
	}
	if value < minVal {
		value = minVal
	} else if value > maxVal {
		value = maxVal
	}
	return value
}

// Abs 计算绝对值
func Abs[T SignedNumber](val T) T {
	return T(math.Abs(float64(val)))
}

// Floor 向下取整
func Floor[T SignedNumber | UnsignedNumber](val T) T {
	return T(math.Floor(float64(val)))
}

// Ceil 向上取整
func Ceil[T SignedNumber | UnsignedNumber](val T) T {
	return T(math.Ceil(float64(val)))
}

// Round 四舍五入取整
func Round[T SignedNumber | UnsignedNumber](val T) T {
	return T(math.Round(float64(val)))
}

// Sum 计算所有元素总和
func Sum[T SignedNumber | UnsignedNumber](items ...T) T {
	var sum T
	for _, item := range items {
		sum += item
	}
	return sum
}

// Average 计算所有元素的平均值
func Average[T SignedNumber | UnsignedNumber](items ...T) T {
	return Sum(items...) / T(len(items))
}
