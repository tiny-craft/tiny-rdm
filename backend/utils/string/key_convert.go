package strutil

import (
	"strconv"
	sliceutil "tinyrdm/backend/utils/slice"
)

// EncodeRedisKey encode the redis key to integer array
// if key contains binary which could not display on ui, convert the key to char array
func EncodeRedisKey(key string) any {
	if ContainsBinary(key) {
		b := []byte(key)
		arr := make([]int, len(b))
		for i, bb := range b {
			arr[i] = int(bb)
		}
		return arr
	}
	return key
}

// DecodeRedisKey decode redis key to readable string
func DecodeRedisKey(key any) string {
	switch key.(type) {
	case string:
		return key.(string)

	case []any:
		arr := key.([]any)
		bytes := sliceutil.Map(arr, func(i int) byte {
			if c, ok := AnyToInt(arr[i]); ok {
				return byte(c)
			}
			return '0'
		})
		return string(bytes)

	case []int:
		arr := key.([]int)
		b := make([]byte, len(arr))
		for i, bb := range arr {
			b[i] = byte(bb)
		}
		return string(b)
	}
	return ""
}

// AnyToInt convert any value to int
func AnyToInt(val any) (int, bool) {
	switch val.(type) {
	case string:
		num, err := strconv.Atoi(val.(string))
		if err != nil {
			return 0, false
		}
		return num, true
	case float64:
		return int(val.(float64)), true
	case float32:
		return int(val.(float32)), true
	case int64:
		return int(val.(int64)), true
	case int32:
		return int(val.(int32)), true
	case int:
		return val.(int), true
	case bool:
		if val.(bool) {
			return 1, true
		} else {
			return 0, true
		}
	}
	return 0, false
}
