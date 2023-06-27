package rand

import (
	"github.com/google/go-cmp/cmp"
	"math/rand"
	"sync"
	"time"
)

// WeightObject 权重单项
type WeightObject[T any] struct {
	Obj    T
	Weight int
}

// WeightRandom 根据权重随机
type WeightRandom[T any] struct {
	WeightObject []WeightObject[T]
	totalWeight  int
	randObj      *rand.Rand
	lk           sync.Mutex
}

func NewWeightRandom[T any]() *WeightRandom[T] {
	return &WeightRandom[T]{
		WeightObject: []WeightObject[T]{},
		totalWeight:  0,
		randObj:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (w *WeightRandom[T]) Add(object T, weight int) {
	weightObj := WeightObject[T]{
		Obj:    object,
		Weight: weight,
	}
	w.AddObject(weightObj)
}

// AddObject 添加单个权重对象
func (w *WeightRandom[T]) AddObject(weightObject WeightObject[T]) {
	if weightObject.Weight <= 0 {
		return
	}

	exists := false
	for i, object := range w.WeightObject {
		if cmp.Equal(weightObject.Obj, object.Obj) {
			// 已经存在, 覆盖权重
			w.subWeight(object.Weight)
			w.WeightObject[i].Weight = weightObject.Weight
			w.addWeight(weightObject.Weight)
			exists = true
			break
		}
	}
	if !exists {
		// 已经存在, 覆盖权重
		w.WeightObject = append(w.WeightObject, weightObject)
		w.addWeight(weightObject.Weight)
	}
}

// AddObjects 添加多个权重对象
func (w *WeightRandom[T]) AddObjects(object []WeightObject[T]) {
	for _, weightObject := range object {
		w.AddObject(weightObject)
	}
}

func (w *WeightRandom[T]) addWeight(weight int) {
	if w.totalWeight < 0 {
		w.totalWeight = 0
	}
	w.totalWeight += weight
}

func (w *WeightRandom[T]) subWeight(weight int) {
	if w.totalWeight-weight < 0 {
		w.totalWeight = 0
	} else {
		w.totalWeight -= weight
	}
}

// Next 通过权重随机到下一个
func (w *WeightRandom[T]) Next() T {
	if w.totalWeight > 0 {
		w.lk.Lock()
		randomWeight := w.randObj.Intn(w.totalWeight)
		w.lk.Unlock()
		weightCount := 0
		for _, object := range w.WeightObject {
			weightCount += object.Weight
			if weightCount > randomWeight {
				return object.Obj
			}
		}
	}
	var noop T
	return noop
}
