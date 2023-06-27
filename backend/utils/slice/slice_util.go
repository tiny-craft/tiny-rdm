package sliceutil

import (
	"sort"
	"strconv"
	"strings"
	. "tinyrdm/backend/utils"
	"tinyrdm/backend/utils/rand"
)

// Get 获取指定索引的值, 如果不存在则返回默认值
func Get[S ~[]T, T any](arr S, index int, defaultVal T) T {
	if index < 0 || index >= len(arr) {
		return defaultVal
	}
	return arr[index]
}

// Remove 删除指定索引的元素
func Remove[S ~[]T, T any](arr S, index int) S {
	return append(arr[:index], arr[index+1:]...)
}

// RemoveIf 移除指定条件的元素
func RemoveIf[S ~[]T, T any](arr S, cond func(T) bool) S {
	l := len(arr)
	if l <= 0 {
		return arr
	}
	for i := l - 1; i >= 0; i-- {
		if cond(arr[i]) {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

// RemoveRange 删除从[from, to]部分元素
func RemoveRange[S ~[]T, T any](arr S, from, to int) S {
	return append(arr[:from], arr[to:]...)
}

// Find 查找指定条件的元素第一个出现位置
func Find[S ~[]T, T any](arr S, matchFunc func(int) bool) (int, bool) {
	total := len(arr)
	for i := 0; i < total; i++ {
		if matchFunc(i) {
			return i, true
		}
	}
	return -1, false
}

// AnyMatch 判断是否有任意元素符合条件
func AnyMatch[S ~[]T, T any](arr S, matchFunc func(int) bool) bool {
	total := len(arr)
	if total > 0 {
		for i := 0; i < total; i++ {
			if matchFunc(i) {
				return true
			}
		}
	}
	return false
}

// AllMatch 判断是否所有元素都符合条件
func AllMatch[S ~[]T, T any](arr S, matchFunc func(int) bool) bool {
	total := len(arr)
	for i := 0; i < total; i++ {
		if !matchFunc(i) {
			return false
		}
	}
	return true
}

// Equals 比较两个切片内容是否完全一致
func Equals[S ~[]T, T comparable](arr1, arr2 S) bool {
	if &arr1 == &arr2 {
		return true
	}

	len1, len2 := len(arr1), len(arr2)
	if len1 != len2 {
		return false
	}
	for i := 0; i < len1; i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

// Contains 判断数组是否包含指定元素
func Contains[S ~[]T, T Hashable](arr S, elem T) bool {
	return AnyMatch(arr, func(idx int) bool {
		return arr[idx] == elem
	})
}

// ContainsAny 判断数组是否包含任意指定元素
func ContainsAny[S ~[]T, T Hashable](arr S, elems ...T) bool {
	for _, elem := range elems {
		if Contains(arr, elem) {
			return true
		}
	}
	return false
}

// ContainsAll 判断数组是否包含所有指定元素
func ContainsAll[S ~[]T, T Hashable](arr S, elems ...T) bool {
	for _, elem := range elems {
		if !Contains(arr, elem) {
			return false
		}
	}
	return true
}

// Filter 筛选出符合指定条件的所有元素
func Filter[S ~[]T, T any](arr S, filterFunc func(int) bool) []T {
	total := len(arr)
	var result []T
	for i := 0; i < total; i++ {
		if filterFunc(i) {
			result = append(result, arr[i])
		}
	}
	return result
}

// Map 数组映射转换
func Map[S ~[]T, T any, R any](arr S, mappingFunc func(int) R) []R {
	total := len(arr)
	result := make([]R, total)
	for i := 0; i < total; i++ {
		result[i] = mappingFunc(i)
	}
	return result
}

// FilterMap 数组过滤和映射转换
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

// ToMap 数组转键值对
func ToMap[S ~[]T, T any, K Hashable, V any](arr S, mappingFunc func(int) (K, V)) map[K]V {
	total := len(arr)
	result := map[K]V{}
	for i := 0; i < total; i++ {
		key, val := mappingFunc(i)
		result[key] = val
	}
	return result
}

// Flat 二维数组扁平化
func Flat[T any](arr [][]T) []T {
	total := len(arr)
	var result []T
	for i := 0; i < total; i++ {
		subTotal := len(arr[i])
		for j := 0; j < subTotal; j++ {
			result = append(result, arr[i][j])
		}
	}
	return result
}

// FlatMap 二维数组扁平化映射
func FlatMap[T any, R any](arr [][]T, mappingFunc func(int, int) R) []R {
	total := len(arr)
	var result []R
	for i := 0; i < total; i++ {
		subTotal := len(arr[i])
		for j := 0; j < subTotal; j++ {
			result = append(result, mappingFunc(i, j))
		}
	}
	return result
}

func FlatValueMap[T Hashable](arr [][]T) []T {
	return FlatMap(arr, func(i, j int) T {
		return arr[i][j]
	})
}

// Reduce 数组累计
func Reduce[S ~[]T, T any, R any](arr S, init R, reduceFunc func(R, T) R) R {
	result := init
	for _, item := range arr {
		result = reduceFunc(result, item)
	}
	return result
}

// Reverse 反转数组(会修改原数组)
func Reverse[S ~[]T, T any](arr S) S {
	total := len(arr)
	for i := 0; i < total/2; i++ {
		arr[i], arr[total-i-1] = arr[total-i-1], arr[i]
	}
	return arr
}

// Join 数组拼接转字符串
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

// JoinString 字符串数组拼接成字符串
func JoinString(arr []string, sep string) string {
	return Join(arr, sep, func(idx int) string {
		return arr[idx]
	})
}

// JoinInt 整形数组拼接转字符串
func JoinInt(arr []int, sep string) string {
	return Join(arr, sep, func(idx int) string {
		return strconv.Itoa(arr[idx])
	})
}

// Unique 数组去重
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

// UniqueEx 数组去重(任意类型)
// @param toKeyFunc 数组元素转为唯一标识字符串函数, 如转为哈希值等
func UniqueEx[S ~[]T, T any](arr S, toKeyFunc func(i int) string) S {
	result := make(S, 0, len(arr))
	keyArr := Map(arr, toKeyFunc)
	uniKeys := map[string]struct{}{}
	var exists bool
	for i, item := range arr {
		if _, exists = uniKeys[keyArr[i]]; !exists {
			uniKeys[keyArr[i]] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Sort 顺序排序(会修改原数组)
func Sort[S ~[]T, T Hashable](arr S) S {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] <= arr[j]
	})
	return arr
}

// SortDesc 倒序排序(会修改原数组)
func SortDesc[S ~[]T, T Hashable](arr S) S {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] > arr[j]
	})
	return arr
}

// Union 返回两个切片共同拥有的元素
func Union[S ~[]T, T Hashable](arr1 S, arr2 S) S {
	hashArr, compArr := arr1, arr2
	if len(arr1) < len(arr2) {
		hashArr, compArr = compArr, hashArr
	}
	hash := map[T]struct{}{}
	for _, item := range hashArr {
		hash[item] = struct{}{}
	}

	uniq := map[T]struct{}{}
	ret := make(S, 0, len(compArr))
	exists := false
	for _, item := range compArr {
		if _, exists = hash[item]; exists {
			if _, exists = uniq[item]; !exists {
				ret = append(ret, item)
				uniq[item] = struct{}{}
			}
		}
	}
	return ret
}

// Exclude 返回不包含的元素
func Exclude[S ~[]T, T Hashable](arr1 S, arr2 S) S {
	diff := make([]T, 0, len(arr1))
	hash := map[T]struct{}{}
	for _, item := range arr2 {
		hash[item] = struct{}{}
	}

	for _, item := range arr1 {
		if _, exists := hash[item]; !exists {
			diff = append(diff, item)
		}
	}
	return diff
}

// PadLeft 左边填充指定数量
func PadLeft[S ~[]T, T any](arr S, val T, count int) S {
	prefix := make(S, count)
	for i := 0; i < count; i++ {
		prefix[i] = val
	}
	arr = append(prefix, arr...)
	return arr
}

// PadRight 右边填充指定数量
func PadRight[S ~[]T, T any](arr S, val T, count int) S {
	for i := 0; i < count; i++ {
		arr = append(arr, val)
	}
	return arr
}

// RemoveLeft 移除左侧相同元素
func RemoveLeft[S ~[]T, T comparable](arr S, val T) S {
	for len(arr) > 0 && arr[0] == val {
		arr = arr[1:]
	}
	return arr
}

// RemoveRight 移除右侧相同元素
func RemoveRight[S ~[]T, T comparable](arr S, val T) S {
	for {
		length := len(arr)
		if length > 0 && arr[length-1] == val {
			arr = arr[:length]
		} else {
			break
		}
	}
	return arr
}

// RandomElem 从切片中随机抽一个
func RandomElem[S ~[]T, T any](arr S) T {
	l := len(arr)
	if l <= 0 {
		var r T
		return r
	}
	return arr[rand.Intn(l)]
}

// RandomElems 从切片中随机抽多个
// 如果切片长度为空, 则返回空切片
func RandomElems[S ~[]T, T any](arr S, count int) []T {
	l := len(arr)
	ret := make([]T, 0, l)
	if l <= 0 {
		return ret
	}

	idxList := rand.IntnCount(l, count)
	for _, idx := range idxList {
		ret = append(ret, arr[idx])
	}
	return ret
}

// RandomUniqElems 从切片中随机抽多个不同的元素
// 如果切片长度为空, 则返回空切片
// 如果所需数量大于切片唯一元素数量, 则返回整个切片
func RandomUniqElems[S ~[]T, T Hashable](arr S, count int) []T {
	if len(arr) <= 0 {
		// 可选列表为空, 返回空切片
		return []T{}
	}
	// 转换为集合
	uniqList := Unique(arr)
	uniqLen := len(uniqList)
	if uniqLen <= count {
		// 可选集合总数<=所需元素数量, 直接返回整个可选集合
		return uniqList
	}

	if count >= uniqLen/2 {
		// 所需唯一元素大于可选集合一半, 随机筛掉(uniqLen-count)个元素
		for i := 0; i < uniqLen-count; i++ {
			uniqList = Remove(uniqList, rand.Intn(uniqLen-i))
		}
		return uniqList
	} else {
		// 所需唯一元素小于可选集合一半, 随机抽取count个元素
		res := make([]T, count)
		var idx int
		for i := 0; i < count; i++ {
			idx = rand.Intn(uniqLen - i)
			res[i] = uniqList[idx]
			uniqList = Remove(uniqList, idx)
		}
		return res
	}
}

// Clone 复制切片
func Clone[S ~[]T, T any](src S) S {
	dest := make(S, len(src))
	copy(dest, src)
	return dest
}

// Count 统计制定条件元素数量
func Count[S ~[]T, T any](arr S, filter func(int) bool) int {
	count := 0
	for i := range arr {
		if filter(i) {
			count += 1
		}
	}
	return count
}

// Group 根据分组函数对数组进行分组汇总
func Group[S ~[]T, T any, K Hashable, R any](arr S, groupFunc func(int) (K, R)) map[K][]R {
	ret := map[K][]R{}
	for i := range arr {
		key, val := groupFunc(i)
		ret[key] = append(ret[key], val)
	}
	return ret
}
