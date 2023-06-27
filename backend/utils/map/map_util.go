package maputil

import (
	. "tinyrdm/backend/utils"
	"tinyrdm/backend/utils/coll"
)

// Get 获取键值对指定键的值, 如果不存在则返回自定默认值
func Get[M ~map[K]V, K Hashable, V any](m M, key K, defaultVal V) V {
	if m != nil {
		if v, exists := m[key]; exists {
			return v
		}
	}
	return defaultVal
}

// ContainsKey 判断指定键是否存在
func ContainsKey[M ~map[K]V, K Hashable, V any](m M, key K) bool {
	if m == nil {
		return false
	}
	_, exists := m[key]
	return exists
}

// MustGet 获取键值对指定键的值, 如果不存在则调用给定的函数进行获取
func MustGet[M ~map[K]V, K Hashable, V any](m M, key K, getFunc func(K) V) V {
	if v, exists := m[key]; exists {
		return v
	}
	if getFunc != nil {
		return getFunc(key)
	}
	var defaultV V
	return defaultV
}

// Keys 获取键值对中所有键
func Keys[M ~map[K]V, K Hashable, V any](m M) []K {
	if len(m) <= 0 {
		return []K{}
	}
	keys := make([]K, len(m))
	index := 0
	for k := range m {
		keys[index] = k
		index += 1
	}
	return keys
}

// KeySet 获取键值对中所有键集合
func KeySet[M ~map[K]V, K Hashable, V any](m M) coll.Set[K] {
	if len(m) <= 0 {
		return coll.NewSet[K]()
	}
	keySet := coll.NewSet[K]()
	for k := range m {
		keySet.Add(k)
	}
	return keySet
}

// Values 获取键值对中所有值
func Values[M ~map[K]V, K Hashable, V any](m M) []V {
	if len(m) <= 0 {
		return []V{}
	}
	values := make([]V, len(m))
	index := 0
	for _, v := range m {
		values[index] = v
		index += 1
	}
	return values
}

// ValueSet 获取键值对中所有值集合
func ValueSet[M ~map[K]V, K Hashable, V Hashable](m M) coll.Set[V] {
	if len(m) <= 0 {
		return coll.NewSet[V]()
	}
	valueSet := coll.NewSet[V]()
	for _, v := range m {
		valueSet.Add(v)
	}
	return valueSet
}

// Fill 填充键值对
func Fill[M ~map[K]V, K Hashable, V any](dest M, src M) M {
	for k, v := range src {
		dest[k] = v
	}
	return dest
}

// Merge 合并键值对, 后续键值对有重复键的元素会覆盖旧元素
func Merge[M ~map[K]V, K Hashable, V any](mapArr ...M) M {
	result := make(M, len(mapArr))
	for _, m := range mapArr {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// DeepMerge 深度递归覆盖src值到dst中
// 将返回新的值
func DeepMerge[M ~map[K]any, K Hashable](src1, src2 M) M {
	out := make(map[K]any, len(src1))
	for k, v := range src1 {
		out[k] = v
	}
	for k, v := range src2 {
		if v1, ok := v.(map[K]any); ok {
			if bv, ok := out[k]; ok {
				if bv1, ok := bv.(map[K]any); ok {
					out[k] = DeepMerge(bv1, v1)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// Omit 根据条件省略指定元素
func Omit[M ~map[K]V, K Hashable, V any](m M, omitFunc func(k K, v V) bool) (M, []K) {
	result := M{}
	var removedKeys []K
	for k, v := range m {
		if !omitFunc(k, v) {
			result[k] = v
		} else {
			removedKeys = append(removedKeys, k)
		}
	}
	return result, removedKeys
}

// OmitKeys 省略指定键的的元素
func OmitKeys[M ~map[K]V, K Hashable, V any](m M, keys ...K) M {
	omitKey := map[K]struct{}{}
	for _, k := range keys {
		omitKey[k] = struct{}{}
	}

	result := M{}
	var exists bool
	for k, v := range m {
		if _, exists = omitKey[k]; !exists {
			result[k] = v
		}
	}
	return result
}

// ContainsAnyKey 是否包含任意键
func ContainsAnyKey[M ~map[K]V, K Hashable, V any](m M, keys ...K) bool {
	var exists bool
	for _, key := range keys {
		if _, exists = m[key]; exists {
			return true
		}
	}
	return false
}

// ContainsAllKey 是否包含所有键
func ContainsAllKey[M ~map[K]V, K Hashable, V any](m M, keys ...K) bool {
	var exists bool
	for _, key := range keys {
		if _, exists = m[key]; !exists {
			return false
		}
	}
	return true
}

// AnyMatch 是否任意元素符合条件
func AnyMatch[M ~map[K]V, K Hashable, V any](m M, matchFunc func(k K, v V) bool) bool {
	for k, v := range m {
		if matchFunc(k, v) {
			return true
		}
	}
	return false
}

// AllMatch 是否所有元素符合条件
func AllMatch[M ~map[K]V, K Hashable, V any](m M, matchFunc func(k K, v V) bool) bool {
	for k, v := range m {
		if !matchFunc(k, v) {
			return false
		}
	}
	return true
}

// Reduce 累计
func Reduce[M ~map[K]V, K Hashable, V any, R any](m M, init R, reduceFunc func(R, K, V) R) R {
	result := init
	for k, v := range m {
		result = reduceFunc(result, k, v)
	}
	return result
}

// ToSlice 键值对转切片
func ToSlice[M ~map[K]V, K Hashable, V any, R any](m M, mapFunc func(k K) R) []R {
	ret := make([]R, 0, len(m))
	for k := range m {
		ret = append(ret, mapFunc(k))
	}
	return ret
}

// Filter 筛选出指定条件的所有元素
func Filter[M ~map[K]V, K Hashable, V any](m M, filterFunc func(k K) bool) M {
	ret := make(M, len(m))
	for k, v := range m {
		if filterFunc(k) {
			ret[k] = v
		}
	}
	return ret
}

// FilterToSlice 键值对筛选并转切片
func FilterToSlice[M ~map[K]V, K Hashable, V any, R any](m M, mapFunc func(k K) (R, bool)) []R {
	ret := make([]R, 0, len(m))
	for k := range m {
		if v, filter := mapFunc(k); filter {
			ret = append(ret, v)
		}
	}
	return ret
}

// FilterKey 筛选出指定条件的所有键
func FilterKey[M ~map[K]V, K Hashable, V any](m M, filterFunc func(k K) bool) []K {
	ret := make([]K, 0, len(m))
	for k := range m {
		if filterFunc(k) {
			ret = append(ret, k)
		}
	}
	return ret
}

// Clone 复制键值对
func Clone[M ~map[K]V, K Hashable, V any](src M) M {
	dest := make(M, len(src))
	for k, v := range src {
		dest[k] = v
	}
	return dest
}

// Reverse 键->值映射翻转为值->键映射(如果重复则覆盖最后的)
func Reverse[M ~map[K]V, K Hashable, V Hashable](src M) map[V]K {
	dest := make(map[V]K, len(src))
	for k, v := range src {
		dest[v] = k
	}
	return dest
}

// ReverseAll 键->值映射翻转为值->键列表映射
func ReverseAll[M ~map[K]V, K Hashable, V Hashable](src M) map[V][]K {
	dest := make(map[V][]K, len(src))
	for k, v := range src {
		dest[v] = append(dest[v], k)
	}
	return dest
}

// RemoveIf 移除指定条件的键
func RemoveIf[M ~map[K]V, K Hashable, V any](src M, cond func(key K) bool) {
	for k := range src {
		if cond(k) {
			delete(src, k)
		}
	}
}
