package coll

// Stack 栈, 先进后出
type Stack[T any] []T

func NewStack[T any](elems ...T) Stack[T] {
	if len(elems) > 0 {
		data := make([]T, len(elems))
		copy(data, elems)
		return data
	} else {
		return Stack[T]{}
	}
}

// Push 顶部添加一个元素
func (s *Stack[T]) Push(elem T) {
	if s == nil {
		panic("queue should not be nil")
	}
	*s = append(*s, elem)
}

// PushN 顶部添加一个元素
func (s *Stack[T]) PushN(elems ...T) {
	if s == nil {
		panic("queue should not be nil")
	}
	if len(elems) <= 0 {
		return
	}
	*s = append(*s, elems...)
}

// Pop 移除并返回顶部元素
func (s *Stack[T]) Pop() T {
	if s == nil {
		panic("queue should not be nil")
	}
	l := len(*s)
	popElem := (*s)[l-1]
	*s = (*s)[:l-1]
	return popElem
}

// PopN 移除并返回顶部多个元素
func (s *Stack[T]) PopN(n int) []T {
	if s == nil {
		panic("queue should not be nil")
	}
	var popElems []T
	if n <= 0 {
		return popElems
	}

	l := len(*s)
	if n >= l {
		popElems = *s
		*s = []T{}
		return *s
	}

	popElems = (*s)[l-n:]
	*s = (*s)[:l-n]

	// 翻转弹出结果
	pl := len(popElems)
	for i := 0; i < pl/2; i++ {
		popElems[i], popElems[pl-i-1] = popElems[pl-i-1], popElems[i]
	}
	return popElems
}

// Clear 移除所有元素
func (s *Stack[T]) Clear() {
	if s == nil {
		panic("queue should not be nil")
	}
	*s = []T{}
}

func (s Stack[T]) IsEmpty() bool {
	return len(s) <= 0
}

func (s Stack[T]) Size() int {
	return len(s)
}
