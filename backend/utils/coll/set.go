package coll

import (
	"fmt"
	json "github.com/bytedance/sonic"
	"sort"
	. "tinyrdm/backend/utils"
	"tinyrdm/backend/utils/rand"
)

type Void struct{}

// Set 集合, 存放不重复的元素
type Set[T Hashable] map[T]Void

// type Set[T Hashable] struct {
//	data map[T]Void
// }

func NewSet[T Hashable](elems ...T) Set[T] {
	if len(elems) > 0 {
		data := make(Set[T], len(elems))
		for _, e := range elems {
			data[e] = Void{}
		}
		return data
	} else {
		return Set[T]{}
	}
}

// Add 添加元素
func (s Set[T]) Add(elem T) bool {
	if s == nil {
		return false
	}
	if _, exists := s[elem]; !exists {
		s[elem] = Void{}
		return true
	}
	return false
}

// AddN 添加多个元素
func (s Set[T]) AddN(elems ...T) int {
	if s == nil {
		return 0
	}
	addCount := 0
	var exists bool
	for _, elem := range elems {
		if _, exists = s[elem]; !exists {
			s[elem] = Void{}
			addCount += 1
		}
	}
	return addCount
}

// Merge 合并其他集合
func (s Set[T]) Merge(other Set[T]) int {
	return s.AddN(other.ToSlice()...)
}

// Contains 判断是否存在指定元素
func (s Set[T]) Contains(elem T) bool {
	if s == nil {
		return false
	}
	_, exists := s[elem]
	return exists
}

// ContainAny 判断是否包含任意元素
func (s Set[T]) ContainAny(elems ...T) bool {
	if s == nil {
		return false
	}
	var exists bool
	for _, elem := range elems {
		if _, exists = s[elem]; exists {
			return true
		}
	}
	return false
}

// Equals 判断两个集合内元素是否一致
func (s Set[T]) Equals(other Set[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for elem := range s {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

// ContainAll 判断是否包含所有元素
func (s Set[T]) ContainAll(elems ...T) bool {
	if s == nil {
		return false
	}
	var exists bool
	for _, elem := range elems {
		if _, exists = s[elem]; !exists {
			return false
		}
	}
	return true
}

// Remove 移除元素
func (s Set[T]) Remove(elem T) bool {
	if s == nil {
		return false
	}
	if _, exists := s[elem]; exists {
		delete(s, elem)
		return true
	}
	return false
}

// RemoveN 移除多个元素
func (s Set[T]) RemoveN(elems ...T) int {
	if s == nil {
		return 0
	}
	var exists bool
	removeCnt := 0
	for _, elem := range elems {
		if _, exists = s[elem]; exists {
			delete(s, elem)
			removeCnt += 1
		}
	}
	return removeCnt
}

// RemoveSub 移除子集
func (s Set[T]) RemoveSub(subSet Set[T]) int {
	if s == nil {
		return 0
	}
	var exists bool
	removeCnt := 0
	for elem := range subSet {
		if _, exists = s[elem]; exists {
			delete(s, elem)
			removeCnt += 1
		}
	}
	return removeCnt
}

// Filter 根据条件筛出符合的元素
func (s Set[T]) Filter(filterFunc func(i T) bool) []T {
	ret := []T{}
	for v := range s {
		if filterFunc(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// RandomElem 随机抽取一个元素
// @param remove 随机出来的元素是否同时从集合中移除
// @return 抽取的元素
// @return 是否抽取成功
func (s Set[T]) RandomElem(remove bool) (T, bool) {
	size := s.Size()
	if size > 0 {
		selIdx := rand.Intn(size)
		idx := 0
		for elem := range s {
			if idx == selIdx {
				if remove {
					delete(s, elem)
				}
				return elem, true
			} else {
				idx++
			}
		}
	}

	var r T
	return r, false
}

// Size 集合长度
func (s Set[T]) Size() int {
	return len(s)
}

// IsEmpty 判断是否为空
func (s Set[T]) IsEmpty() bool {
	return len(s) <= 0
}

// Clear 清空集合
func (s Set[T]) Clear() {
	for elem := range s {
		delete(s, elem)
	}
}

// ToSlice 转为切片
func (s Set[T]) ToSlice() []T {
	size := len(s)
	if size <= 0 {
		return []T{}
	}

	ret := make([]T, 0, size)
	for elem := range s {
		ret = append(ret, elem)
	}
	return ret
}

// ToSortedSlice 转为排序好的切片
func (s Set[T]) ToSortedSlice(sortFunc func(v1, v2 T) bool) []T {
	list := s.ToSlice()
	sort.Slice(list, func(i, j int) bool {
		return sortFunc(list[i], list[j])
	})
	return list
}

// Each 遍历检索每个元素
func (s Set[T]) Each(eachFunc func(T)) {
	if len(s) <= 0 {
		return
	}
	for elem := range s {
		eachFunc(elem)
	}
}

// Clone 克隆
func (s Set[T]) Clone() Set[T] {
	if s == nil {
		return nil
	}

	other := NewSet[T]()
	for elem := range s {
		other[elem] = Void{}
	}
	return other
}

func (s Set[T]) String() string {
	arr := s.ToSlice()
	return fmt.Sprintf("%v", arr)
}

// MarshalJSON to output non base64 encoded []byte
func (s Set[T]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	t := s.ToSlice()
	return json.Marshal(t)
}

// UnmarshalJSON to deserialize []byte
func (s *Set[T]) UnmarshalJSON(b []byte) error {
	t := []T{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		*s = NewSet[T]()
	} else {
		*s = NewSet[T](t...)
	}
	return nil
}

// GormDataType gorm common data type
func (s Set[T]) GormDataType() string {
	return "json"
}
