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

// 将interface{}切片转换为字符串切片，如果可能
func interfaceSliceToStringSlice(slice []interface{}) ([]string, bool) {
	strSlice := make([]string, 0, len(slice))
	for _, item := range slice {
		if str, ok := item.(string); ok {
			strSlice = append(strSlice, str)
		} else {
			return nil, false // 如果存在非字符串元素，则返回错误
		}
	}
	return strSlice, true
}

func interfaceSliceToSortSlice(slice []interface{}) ([]interface{}, bool) {
	strSlice := make([]string, 0, len(slice))
	for _, item := range slice {
		if b, err := json.Marshal(item); err == nil {
			strSlice = append(strSlice, string(b))
		} else {
			return nil, false // 如果存在非字符串元素，则返回错误
		}
	}
	sort.Strings(strSlice)
	sortSlice := make([]interface{}, 0, len(strSlice))
	for _, item := range strSlice {
		var itemV interface{}
		if err := json.Unmarshal([]byte(item), &itemV); err == nil {
			sortSlice = append(sortSlice, itemV)
		} else {
			return nil, false // 如果存在非字符串元素，则返回错误
		}
	}
	return sortSlice, true
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
			// 尝试将切片转换为字符串切片
			if strSlice, ok := v.([]string); ok {
				sort.Strings(strSlice)
				result[k] = strSlice
			} else if slice, ok := v.([]interface{}); ok {
				// 如果切片元素是interface{}，尝试转换为字符串切片
				if strSlice, ok := interfaceSliceToStringSlice(slice); ok {
					sort.Strings(strSlice)
					result[k] = strSlice
				} else if sortSlice, ok := interfaceSliceToSortSlice(slice); ok {
					result[k] = sortSlice
				} else {
					// 对于非字符串切片，保持原样或进行其他处理
					result[k] = slice
				}
			}
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

	result := make(map[string]interface{})
	for k, v := range m {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			// 尝试将切片转换为字符串切片
			if strSlice, ok := v.([]string); ok {
				sort.Strings(strSlice)
				result[k] = strSlice
			} else if slice, ok := v.([]interface{}); ok {
				// 如果切片元素是interface{}，尝试转换为字符串切片
				if strSlice, ok := interfaceSliceToStringSlice(slice); ok {
					sort.Strings(strSlice)
					result[k] = strSlice
				} else if sortSlice, ok := interfaceSliceToSortSlice(slice); ok {
					result[k] = sortSlice
				} else {
					// 对于非字符串切片，保持原样或进行其他处理
					result[k] = slice
				}
			}
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
	bytes, _ := json.Marshal(result)
	return bytes
}
