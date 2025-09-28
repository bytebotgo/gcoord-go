// Package gcoord 提供坐标系转换，参考 gcoord-js。
//
// 支持的坐标系：
//   - WGS84: 世界大地坐标系，GPS 原始坐标
//   - GCJ02: 国测局坐标系，中国标准，火星坐标
//   - BD09: 百度坐标系，百度地图使用
//   - BD09MC: 百度墨卡托投影坐标系
//   - EPSG3857: Web 墨卡托投影坐标系，Google Maps 等使用
//
// 精度说明：
//   - 经纬度转换精度：约 1e-5 度（约 1 米）
//   - 投影坐标转换精度：约 1 米
//   - WGS84↔GCJ02 在中国境外无偏移
//
// 用法示例：
//
//	// 基础坐标转换
//	p, _ := gcoord.Transform(gcoord.Position{116.397, 39.908}, gcoord.WGS84, gcoord.GCJ02)
//	// p => GCJ02 坐标
//
//	// GeoJSON 转换
//	feature := map[string]any{
//	    "type": "Feature",
//	    "geometry": map[string]any{
//	        "type": "Point",
//	        "coordinates": []any{116.397, 39.908},
//	    },
//	}
//	converted, _ := gcoord.Transform(feature, gcoord.WGS84, gcoord.GCJ02)
//
//	// JSON 字符串转换
//	jsonStr := `{"type":"Point","coordinates":[116.397,39.908]}`
//	result, _ := gcoord.Transform(jsonStr, gcoord.WGS84, gcoord.BD09)
package gcoord
