package coll

// Queue 队列, 先进先出
type Queue[T any] []T

func NewQueue[T any](elems ...T) Queue[T] {
	if len(elems) > 0 {
		data := make([]T, len(elems))
		copy(data, elems)
		return data
	} else {
		return Queue[T]{}
	}
}

// Push 尾部插入元素
func (q *Queue[T]) Push(elem T) {
	if q == nil {
		return
	}
	*q = append(*q, elem)
}

// PushN 尾部插入多个元素
func (q *Queue[T]) PushN(elems ...T) {
	if q == nil {
		return
	}
	if len(elems) <= 0 {
		return
	}
	*q = append(*q, elems...)
}

// Pop 移除并返回头部元素
func (q *Queue[T]) Pop() (T, bool) {
	var elem T
	if q == nil || len(*q) <= 0 {
		return elem, false
	}
	elem = (*q)[0]
	*q = (*q)[1:]
	return elem, true
}

func (q *Queue[T]) PopN(n int) []T {
	if q == nil {
		return []T{}
	}

	var popElems []T
	if n <= 0 {
		return []T{}
	}

	l := len(*q)
	if n >= l {
		popElems = *q
		*q = []T{}
		return *q
	}

	popElems = (*q)[:n]
	*q = (*q)[n:]
	return popElems
}

// Clear 移除所有元素
func (q *Queue[T]) Clear() {
	if q == nil {
		return
	}
	*q = []T{}
}

func (q Queue[T]) IsEmpty() bool {
	return len(q) <= 0
}

func (q Queue[T]) Size() int {
	return len(q)
}
