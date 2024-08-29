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

func interfaceSliceToSortSlice(slice interface{}) (interface{}, bool) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return nil, false
	}

	switch v := slice.(type) {
	case []string:
		SortSlice(v)
	case []int:
		SortSlice(v)
	case []int8:
		SortSlice(v)
	case []int16:
		SortSlice(v)
	case []int32:
		SortSlice(v)
	case []int64:
		SortSlice(v)
	case []uint:
		SortSlice(v)
	case []uint8:
		SortSlice(v)
	case []uint16:
		SortSlice(v)
	case []uint32:
		SortSlice(v)
	case []uint64:
		SortSlice(v)
	case []float32:
		SortSlice(v)
	case []float64:
		SortSlice(v)
	case []uintptr:
		SortSlice(v)
	case []map[string]interface{}:
		sSlice := make([]string, 0)
		for _, v := range v {
			vb, _ := json.Marshal(SortMapValues(v))
			sSlice = append(sSlice, string(vb))
		}
		sort.Strings(sSlice)
		var ret = make([]map[string]interface{}, len(sSlice))
		for i, s := range sSlice {
			json.Unmarshal([]byte(s), &ret[i])
		}
		slice = ret
	}

	return slice, true
}

// 对map的值进行递归处理，如果值是数组或切片，则排序
func SortMapValues(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if v == nil {
			result[k] = ""
			continue
		}

		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			result[k], _ = interfaceSliceToSortSlice(v)
		case reflect.Map:
			// 递归处理嵌套的map
			nestedMap := v.(map[string]interface{})
			sortedMap := SortMapValues(nestedMap)
			result[k] = sortedMap
		default:
			// 其他类型，直接赋值
			result[k] = v
		}
	}
	return result
}

// 对map的值进行递归处理，如果值是数组或切片，则排序
func SortMapValuesByBytes(jsonBytes []byte) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &m); err != nil {
		return nil
	}

	result := SortMapValues(m)
	bytes, _ := json.Marshal(result)

	return bytes
}
