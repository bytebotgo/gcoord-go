package gcoord

import (
	"encoding/json"
	"math"
)

// transformAny 递归遍历 GeoJSON，原地转换 coordinates
func transformAny(obj any, conv Converter) any {
	switch t := obj.(type) {
	case map[string]any:
		// FeatureCollection / Feature / Geometry
		if ty, ok := t["type"].(string); ok {
			switch ty {
			case "FeatureCollection":
				if arr, ok := t["features"].([]any); ok {
					for i := range arr {
						arr[i] = transformAny(arr[i], conv)
					}
					t["features"] = arr
				}
			case "Feature":
				if g, ok := t["geometry"].(map[string]any); ok {
					t["geometry"] = transformAny(g, conv)
				}
			case "GeometryCollection":
				if geoms, ok := t["geometries"].([]any); ok {
					for i := range geoms {
						if gm, ok := geoms[i].(map[string]any); ok {
							geoms[i] = transformAny(gm, conv)
						}
					}
					t["geometries"] = geoms
				}
			default:
				// Geometry
				if coords, ok := t["coordinates"]; ok {
					t["coordinates"] = transformCoords(coords, conv)
				}
			}
		}
		return t
	case []any:
		for i := range t {
			t[i] = transformAny(t[i], conv)
		}
		return t
	default:
		return obj
	}
}

// transformCoords 转换坐标数组
func transformCoords(coords any, conv Converter) any {
	switch c := coords.(type) {
	case []any:
		// 可能是 [x,y] 或 多维数组
		if len(c) >= 2 {
			// 判断是否为数值型坐标（允许超过2维，保留额外维度）
			x := toFloat(c[0])
			y := toFloat(c[1])
			if !math.IsNaN(x) && !math.IsNaN(y) && (len(c) == 2 || isAllNumbers(c)) {
				res := conv(Position{x, y})
				out := make([]any, len(c))
				out[0] = res[0]
				out[1] = res[1]
				if len(c) > 2 {
					copy(out[2:], c[2:])
				}
				return out
			}
		}
		for i := range c {
			c[i] = transformCoords(c[i], conv)
		}
		return c
	default:
		return coords
	}
}

// isAllNumbers 检查数组中的所有元素是否都是数字
func isAllNumbers(arr []any) bool {
	for _, v := range arr {
		if _, ok := v.(float64); !ok {
			return false
		}
	}
	return true
}

// toFloat 将任意类型转换为 float64
func toFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case json.Number:
		f, _ := n.Float64()
		return f
	default:
		return math.NaN()
	}
}
