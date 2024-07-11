package pogo

import "strconv"

type Param map[string]any

// Exists key是否存在，返回bool值
func (p Param) Exists(key string) bool {
	_, exists := p[key]
	return exists
}

// Bool 获取bool值
func (p Param) Bool(key string, defaultVal bool) bool {
	if value, ok := p[key].(bool); ok {
		return value
	}

	return defaultVal
}

// String 获取string值
func (p Param) String(key, defaultVal string) string {
	if value, ok := p[key].(string); ok {
		return value
	}

	return defaultVal
}

// StringSlice 获取[]string
func (p Param) StringSlice(key string, defaultVal []string) []string {
	value, ok := p[key].([]any)
	if !ok {
		return defaultVal
	}

	// 空切片返回默认数据
	if len(value) < 1 {
		return defaultVal
	}

	var stringList = make([]string, 0, len(value))
	for index := 0; index < len(value); index++ {
		// 非[]string，返回默认数据
		val, ok := value[index].(string)
		if !ok {
			return defaultVal
		}
		stringList = append(stringList, val)
	}

	return stringList
}

// Int 获取Int
func (p Param) Int(key string, defaultVal int) int {
	return int(p.Int64(key, int64(defaultVal)))
}

// IntSlice 获取[]int
func (p Param) IntSlice(key string, defaultVal []int) []int {
	value, ok := p[key].([]any)
	if !ok {
		return defaultVal
	}

	if len(value) < 1 {
		return defaultVal
	}

	var intList = make([]int, 0, len(value))
	for index := 0; index < len(value); index++ {
		val, ok := value[index].(float64)
		if !ok {
			return defaultVal
		}
		intList = append(intList, int(val))
	}

	return intList
}

// Int64 获取Int64
func (p Param) Int64(key string, defaultVal int64) int64 {
	if value, ok := p[key].(float64); ok {
		return int64(value)
	}

	return defaultVal
}

// Int64Slice 获取[]int64
func (p Param) Int64Slice(key string, defaultVal []int64) []int64 {
	value, ok := p[key].([]any)
	if !ok {
		return defaultVal
	}

	if len(value) < 1 {
		return defaultVal
	}

	var int64List = make([]int64, 0, len(value))
	for index := 0; index < len(value); index++ {
		val, ok := value[index].(float64)
		if !ok {
			return defaultVal
		}
		int64List = append(int64List, int64(val))
	}

	return int64List
}

// Uint32 Uint32
func (p Param) Uint32(key string, defaultVal uint32) uint32 {
	if value, ok := p[key].(float64); ok {
		return uint32(value)
	}

	return defaultVal
}

// Uint32Slice 获取[]uint32
func (p Param) Uint32Slice(key string, defaultVal []uint32) []uint32 {
	value, ok := p[key].([]any)
	if !ok {
		return defaultVal
	}

	if len(value) < 1 {
		return defaultVal
	}

	var uint32List = make([]uint32, 0, len(value))
	for index := 0; index < len(value); index++ {
		val, ok := value[index].(float64)
		if !ok {
			return defaultVal
		}
		uint32List = append(uint32List, uint32(val))
	}

	return uint32List
}

// Float64 获取float64
func (p Param) Float64(key string, defaultVal float64) float64 {
	if value, ok := p[key].(float64); ok {
		return value
	}
	return defaultVal
}

// Uint64 获取Uint64
func (p Param) Uint64(key string, defaultVal uint64) uint64 {
	return uint64(p.Int64(key, int64(defaultVal)))
}

// IndexStringMap 获取map[int]string
func (p Param) IndexStringMap(key string, defaultVal map[int]string) map[int]string {
	value, exists := p[key].(map[string]any)
	if !exists {
		return defaultVal
	}

	val := make(map[int]string, len(value))
	if len(value) == 0 {
		return val
	}

	for k, v := range value {
		index, err := strconv.Atoi(k)
		mV, ok := v.(string)
		if err != nil || !ok {
			continue
		}

		val[index] = mV
	}

	return val
}
