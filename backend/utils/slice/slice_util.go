package sliceutil

import (
	"strings"
	. "tinyrdm/backend/utils"
)

// Map map items to new array
func Map[S ~[]T, T any, R any](arr S, mappingFunc func(int) R) []R {
	total := len(arr)
	result := make([]R, total)
	for i := 0; i < total; i++ {
		result[i] = mappingFunc(i)
	}
	return result
}

// FilterMap filter and map items to new array
func FilterMap[S ~[]T, T any, R any](arr S, mappingFunc func(int) (R, bool)) []R {
	total := len(arr)
	result := make([]R, 0, total)
	var filter bool
	var mapItem R
	for i := 0; i < total; i++ {
		if mapItem, filter = mappingFunc(i); filter {
			result = append(result, mapItem)
		}
	}
	return result
}

// Join join any array to a single string by custom function
func Join[S ~[]T, T any](arr S, sep string, toStringFunc func(int) string) string {
	total := len(arr)
	if total <= 0 {
		return ""
	}
	if total == 1 {
		return toStringFunc(0)
	}

	sb := strings.Builder{}
	for i := 0; i < total; i++ {
		if i != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(toStringFunc(i))
	}
	return sb.String()
}

// JoinString join string array to a single string
func JoinString(arr []string, sep string) string {
	return Join(arr, sep, func(idx int) string {
		return arr[idx]
	})
}

// Unique filter unique item
func Unique[S ~[]T, T Hashable](arr S) S {
	result := make(S, 0, len(arr))
	uniKeys := map[T]struct{}{}
	var exists bool
	for _, item := range arr {
		if _, exists = uniKeys[item]; !exists {
			uniKeys[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
