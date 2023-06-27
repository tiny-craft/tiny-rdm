package rand

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

// 随机对象缓存池(解决自带随机函数全局抢锁问题)
var randObjectPool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}
var lowerChar = []rune("abcdefghijklmnopqrstuvwxyz") // strings.Split("abcdefghijklmnopqrstuvwxyz", "")
var upperChar = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ") // strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
var numberChar = []rune("0123456789")                // strings.Split("0123456789", "")
var numberAndChar = append(lowerChar, numberChar...)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Intn[T ~int](n T) T {
	r := randObjectPool.Get()
	res := r.(*rand.Rand).Intn(int(n))
	randObjectPool.Put(r)
	return T(res)
}

func IntnCount[T ~int](n T, count int) []T {
	res := make([]T, count)
	r := randObjectPool.Get()
	for i := 0; i < count; i++ {
		res[i] = T(r.(*rand.Rand).Intn(int(n)))
	}
	randObjectPool.Put(r)
	return res
}

// Int31n 生成<n的32位整形
func Int31n[T ~int32](n T) T {
	r := randObjectPool.Get()
	res := r.(*rand.Rand).Int31n(int32(n))
	randObjectPool.Put(r)
	return T(res)
}

// Int63n 生成小于n的64位整形
func Int63n[T ~int64](n T) T {
	r := randObjectPool.Get()
	res := r.(*rand.Rand).Int63n(int64(n))
	randObjectPool.Put(r)
	return T(res)
}

// RangeInt 获取范围内的随机整数[min, max)
func RangeInt[T ~int](min, max T) T {
	if min > max {
		min, max = max, min
	}
	return Intn(max-min) + min
}

// RangeString 生成随机字符串
func RangeString(charSet []rune, n int) string {
	r := randObjectPool.Get()

	res := strings.Builder{}
	size := len(charSet)
	for i := 0; i < n; i++ {
		res.WriteRune(charSet[r.(*rand.Rand).Intn(size)])
	}
	randObjectPool.Put(r)
	return res.String()
}

// LowerString 生成随机指定长度小写字母
func LowerString(n int) string {
	return RangeString(lowerChar, n)
}

// UpperString 生成随机指定长度大写字母
func UpperString(n int) string {
	return RangeString(upperChar, n)
}

// NumberString 生成随机指定长度数字字符串
func NumberString(n int) string {
	return RangeString(numberChar, n)
}

// CharNumberString 生成随机指定长度小写字母和数字
func CharNumberString(n int) string {
	return RangeString(numberAndChar, n)
}

// Shuffle 执行指定次数打乱
func Shuffle(n int, swap func(i, j int)) {
	r := randObjectPool.Get()
	r.(*rand.Rand).Shuffle(n, swap)
}
