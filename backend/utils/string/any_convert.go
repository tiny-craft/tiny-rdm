package strutil

import (
	"encoding/json"
	"strconv"
	sliceutil "tinyrdm/backend/utils/slice"
)

func AnyToString(value interface{}) (s string) {
	if value == nil {
		return
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		s = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		s = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		s = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		s = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		s = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		s = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		s = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		s = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		s = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		s = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		s = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		s = strconv.FormatUint(it, 10)
	case string:
		s = value.(string)
	case bool:
		val, _ := value.(bool)
		if val {
			s = "True"
		} else {
			s = "False"
		}
	case []byte:
		s = string(value.([]byte))
	case []string:
		ss := value.([]string)
		anyStr := sliceutil.Map(ss, func(i int) string {
			str := AnyToString(ss[i])
			return strconv.Itoa(i+1) + ") \"" + str + "\""
		})
		s = sliceutil.JoinString(anyStr, "\r\n")
	case []any:
		as := value.([]any)
		anyItems := sliceutil.Map(as, func(i int) string {
			str := AnyToString(as[i])
			return strconv.Itoa(i+1) + ") \"" + str + "\""
		})
		s = sliceutil.JoinString(anyItems, "\r\n")
	default:
		b, _ := json.Marshal(value)
		s = string(b)
	}

	return
}

//func AnyToHex(val any) (string, bool) {
//	var src string
//	switch val.(type) {
//	case string:
//		src = val.(string)
//	case []byte:
//		src = string(val.([]byte))
//	}
//
//	if len(src) <= 0 {
//		return "", false
//	}
//
//	var output strings.Builder
//	for i := range src {
//		if !utf8.ValidString(src[i : i+1]) {
//			output.WriteString(fmt.Sprintf("\\x%02x", src[i:i+1]))
//		} else {
//			output.WriteString(src[i : i+1])
//		}
//	}
//
//	return output.String(), true
//}
