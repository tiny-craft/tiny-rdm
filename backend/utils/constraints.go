package utils

type Hashable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

type SignedNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type UnsignedNumber interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
