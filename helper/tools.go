package helper

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
	"time"
)

// Uint64List 排序
type Uint64List []uint64

func (my64 Uint64List) Len() int           { return len(my64) }
func (my64 Uint64List) Swap(i, j int)      { my64[i], my64[j] = my64[j], my64[i] }
func (my64 Uint64List) Less(i, j int) bool { return my64[i] < my64[j] }

// Comparable 定义可比较的类型
type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// SortSlice 排序
func SortSlice[T Comparable](slice []T) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
}

// Md5 生成md5
func Md5(str []byte) string {
	hash := md5.Sum(str)
	return hex.EncodeToString(hash[:])
}

// DiffNano 时间差，纳秒
func DiffNano(startTime time.Time) (diff int64) {
	diff = int64(time.Since(startTime))
	return
}

// InArrayStr 判断字符串是否在数组内
func InArrayStr(str string, arr []string) (inArray bool) {
	for _, s := range arr {
		if s == str {
			inArray = true
			break
		}
	}
	return
}

func interfaceSliceToSortSlice(slice []interface{}) (interface{}, bool) {
	if len(slice) == 0 || slice[0] == nil {
		return slice, true
	}

	llen := len(slice)
	var ret interface{}

	switch slice[0].(type) {
	case string:
		typedSlice := make([]string, llen)
		for i, v := range slice {
			typedSlice[i] = v.(string)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case int:
		typedSlice := make([]int, llen)
		for i, v := range slice {
			typedSlice[i] = v.(int)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case int8:
		typedSlice := make([]int8, llen)
		for i, v := range slice {
			typedSlice[i] = v.(int8)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case int16:
		typedSlice := make([]int16, llen)
		for i, v := range slice {
			typedSlice[i] = v.(int16)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case int32:
		typedSlice := make([]int32, llen)
		for i, v := range slice {
			typedSlice[i] = v.(int32)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case int64:
		typedSlice := make([]int64, llen)
		for i, v := range slice {
			typedSlice[i] = v.(int64)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case uint:
		typedSlice := make([]uint, llen)
		for i, v := range slice {
			typedSlice[i] = v.(uint)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case uint8:
		typedSlice := make([]uint8, llen)
		for i, v := range slice {
			typedSlice[i] = v.(uint8)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case uint16:
		typedSlice := make([]uint16, llen)
		for i, v := range slice {
			typedSlice[i] = v.(uint16)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case uint32:
		typedSlice := make([]uint32, llen)
		for i, v := range slice {
			typedSlice[i] = v.(uint32)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case uint64:
		typedSlice := make([]uint64, llen)
		for i, v := range slice {
			typedSlice[i] = v.(uint64)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case float32:
		typedSlice := make([]float32, llen)
		for i, v := range slice {
			typedSlice[i] = v.(float32)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case float64:
		typedSlice := make([]float64, llen)
		for i, v := range slice {
			typedSlice[i] = v.(float64)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case uintptr:
		typedSlice := make([]uintptr, llen)
		for i, v := range slice {
			typedSlice[i] = v.(uintptr)
		}
		SortSlice(typedSlice)
		ret = typedSlice
	case map[string]interface{}:
		sSlice := make([]string, 0, len(slice))
		for i, v := range slice {
			vv := sortMapValues(v.(map[string]interface{}))
			slice[i] = vv

			vb, _ := json.Marshal(vv)
			sSlice = append(sSlice, string(vb))
		}
		sort.Strings(sSlice)
		for i, s := range sSlice {
			json.Unmarshal([]byte(s), &slice[i])
		}
		ret = slice
	default:
		ret = slice
	}
	return ret, true
}

// sortMapValues 对map的值进行递归处理，如果值是数组或切片，则排序（仅接收从json反解析后的map，自定义的map这里会报错）
func sortMapValues(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if v == nil {
			result[k] = ""
			continue
		}

		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			result[k], _ = interfaceSliceToSortSlice(v.([]interface{}))
		case reflect.Map:
			// 递归处理嵌套的map
			nestedMap := v.(map[string]interface{})
			sortedMap := sortMapValues(nestedMap)
			result[k] = sortedMap
		default:
			// 其他类型，直接赋值
			result[k] = v
		}
	}
	return result
}

// SortMapValuesByBytes 对map的值进行递归处理，如果值是数组或切片，则排序
func SortMapValuesByBytes(jsonBytes []byte) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &m); err != nil {
		return nil
	}

	result := sortMapValues(m)
	bytes, _ := json.Marshal(result)

	return bytes
}

// SortMapValuesByMap 对map的值进行递归处理，如果值是数组或切片，则排序 (针对定义好的map需要重新编码)
func SortMapValuesByMap(m map[string]interface{}) map[string]interface{} {
	body, _ := json.Marshal(m)
	var mm map[string]interface{}
	json.Unmarshal(body, &mm)

	return sortMapValues(mm)
}
