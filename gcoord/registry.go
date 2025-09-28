package gcoord

import (
	"sync"
)

// Converter 将一个坐标转换为另一个坐标
type Converter func(Position) Position

// crsMap 记录从某 CRS 到其他 CRS 的转换函数
var crsMap = map[CRSTypes]map[CRSTypes]Converter{}

// 预计算的转换器缓存
var (
	converterCache = make(map[string]Converter)
	cacheMutex     sync.RWMutex
)

func init() {
	// 注册各 CRS 的转换函数
	crsMap[WGS84] = map[CRSTypes]Converter{
		GCJ02:    WGS84ToGCJ02,
		BD09:     compose(GCJ02ToBD09, WGS84ToGCJ02),
		BD09MC:   compose(BD09toBD09MC, GCJ02ToBD09, WGS84ToGCJ02),
		EPSG3857: WGS84ToEPSG3857,
	}
	crsMap[GCJ02] = map[CRSTypes]Converter{
		WGS84:    GCJ02ToWGS84,
		BD09:     GCJ02ToBD09,
		BD09MC:   compose(BD09toBD09MC, GCJ02ToBD09),
		EPSG3857: compose(WGS84ToEPSG3857, GCJ02ToWGS84),
	}
	crsMap[BD09] = map[CRSTypes]Converter{
		WGS84:    compose(GCJ02ToWGS84, BD09ToGCJ02),
		GCJ02:    BD09ToGCJ02,
		EPSG3857: compose(WGS84ToEPSG3857, GCJ02ToWGS84, BD09ToGCJ02),
		BD09MC:   BD09toBD09MC,
	}
	crsMap[EPSG3857] = map[CRSTypes]Converter{
		WGS84:  EPSG3857ToWGS84,
		GCJ02:  compose(WGS84ToGCJ02, EPSG3857ToWGS84),
		BD09:   compose(GCJ02ToBD09, WGS84ToGCJ02, EPSG3857ToWGS84),
		BD09MC: compose(BD09toBD09MC, GCJ02ToBD09, WGS84ToGCJ02, EPSG3857ToWGS84),
	}
	crsMap[BD09MC] = map[CRSTypes]Converter{
		WGS84:    compose(GCJ02ToWGS84, BD09ToGCJ02, BD09MCtoBD09),
		GCJ02:    compose(BD09ToGCJ02, BD09MCtoBD09),
		EPSG3857: compose(WGS84ToEPSG3857, GCJ02ToWGS84, BD09ToGCJ02, BD09MCtoBD09),
		BD09:     BD09MCtoBD09,
	}
}

// compose 将多个 Converter 组合为一个，从右到左执行
func compose(funcs ...Converter) Converter {
	return func(p Position) Position {
		res := p
		for i := len(funcs) - 1; i >= 0; i-- {
			res = funcs[i](res)
		}
		return res
	}
}

// getConverter 获取或创建转换器，支持缓存
func getConverter(from, to CRSTypes) Converter {
	if from == to {
		return func(p Position) Position { return p }
	}

	key := string(from) + "->" + string(to)

	cacheMutex.RLock()
	if converter, exists := converterCache[key]; exists {
		cacheMutex.RUnlock()
		return converter
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 双重检查
	if converter, exists := converterCache[key]; exists {
		return converter
	}

	// 创建转换器
	converter := buildConverter(from, to)
	if converter != nil {
		converterCache[key] = converter
	}
	return converter
}

// buildConverter 构建转换器
func buildConverter(from, to CRSTypes) Converter {
	fromMap, ok := crsMap[from]
	if !ok {
		return nil
	}

	conv, ok := fromMap[to]
	if !ok || conv == nil {
		return nil
	}

	return conv
}

// ClearCache 清空转换器缓存（主要用于测试）
func ClearCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	converterCache = make(map[string]Converter)
}
