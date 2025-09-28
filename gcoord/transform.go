package gcoord

import (
	"encoding/json"
	"fmt"
)

// 转换器注册和组合逻辑已移至 registry.go

// Transform 将输入从 crsFrom 转换到 crsTo。
//
// 支持的输入类型：
//   - Position: []float64{lon, lat} 或 []float64{lon, lat, ...}
//   - string: JSON 字符串（Position 或 GeoJSON 对象）
//   - map[string]any: 任意 GeoJSON 对象（Point/LineString/Polygon/Feature/FeatureCollection/...）
//   - []any: 坐标数组
//
// 转换精度：
//   - 经纬度转换：约 1e-5 度（约 1 米）
//   - 投影坐标转换：约 1 米
//
// 示例：
//
//	p, _ := Transform(Position{116.397, 39.908}, WGS84, GCJ02)
//	feature, _ := Transform(geoJSON, WGS84, BD09)
//	result, _ := Transform(`{"type":"Point","coordinates":[116.397,39.908]}`, WGS84, EPSG3857)
func Transform[T any](input T, crsFrom, crsTo CRSTypes) (T, error) {
	var zero T

	// 验证输入参数
	if err := validateCRS(crsFrom); err != nil {
		return zero, err
	}
	if err := validateCRS(crsTo); err != nil {
		return zero, err
	}

	if crsFrom == crsTo {
		return input, nil
	}

	conv := getConverter(crsFrom, crsTo)
	if conv == nil {
		return zero, fmt.Errorf("无效的目标坐标系: %s", crsTo)
	}

	// 尝试类型分支
	switch v := any(input).(type) {
	case string:
		var obj any
		if err := json.Unmarshal([]byte(v), &obj); err != nil {
			return zero, ErrJSONParseFailed(err)
		}
		out := transformAny(obj, conv)
		b, _ := json.Marshal(out)
		return any(string(b)).(T), nil
	case Position:
		if err := validatePosition(v); err != nil {
			return zero, err
		}
		v = ensureNumberSlice(v)
		return any(conv(v)).(T), nil
	case []float64:
		if err := validatePosition(Position(v)); err != nil {
			return zero, err
		}
		return any(conv(Position(v))).(T), nil
	default:
		out := transformAny(v, conv)
		return any(out).(T), nil
	}
}

// 工具函数已移至 utils.go
// GeoJSON 处理函数已移至 geojson.go
