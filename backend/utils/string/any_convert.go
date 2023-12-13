package strutil

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	sliceutil "tinyrdm/backend/utils/slice"
)

func AnyToString(value interface{}, prefix string, layer int) (s string) {
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
		if layer > 0 {
			s = "\"" + value.(string) + "\""
		} else {
			s = value.(string)
		}
	case bool:
		val, _ := value.(bool)
		if val {
			s = "True"
		} else {
			s = "False"
		}
	case []byte:
		s = prefix + string(value.([]byte))
	case []string:
		ss := value.([]string)
		anyStr := sliceutil.Map(ss, func(i int) string {
			str := AnyToString(ss[i], prefix, layer+1)
			return prefix + strconv.Itoa(i+1) + ") " + str
		})
		s = prefix + sliceutil.JoinString(anyStr, "\r\n")
	case []any:
		as := value.([]any)
		anyItems := sliceutil.Map(as, func(i int) string {
			str := AnyToString(as[i], prefix, layer+1)
			return prefix + strconv.Itoa(i+1) + ") " + str
		})
		s = sliceutil.JoinString(anyItems, "\r\n")
	case map[any]any:
		am := value.(map[any]any)
		var items []string
		index := 0
		for k, v := range am {
			kk := prefix + strconv.Itoa(index+1) + ") " + AnyToString(k, prefix, layer+1)
			vv := prefix + strconv.Itoa(index+2) + ") " + AnyToString(v, "\t", layer+1)
			if layer > 0 {
				indent := layer
				if index == 0 {
					indent -= 1
				}
				for i := 0; i < indent; i++ {
					vv = "  " + vv
				}
			}
			index += 2
			items = append(items, kk, vv)
		}
		s = sliceutil.JoinString(items, "\r\n")
	default:
		b, _ := json.Marshal(value)
		s = prefix + string(b)
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

func SplitCmd(cmd string) []string {
	re := regexp.MustCompile(`'[^']+'|"[^"]+"|\S+`)
	args := re.FindAllString(cmd, -1)
	return sliceutil.FilterMap(args, func(i int) (string, bool) {
		arg := strings.Trim(args[i], "\"")
		arg = strings.Trim(arg, "'")
		if len(arg) <= 0 {
			return "", false
		}
		return arg, true
	})
}
