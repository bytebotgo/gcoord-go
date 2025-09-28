# gcoord-go

Go 版本的坐标系转换库，参考 [gcoord-js](https://github.com/hujiulong/gcoord) 实现。

## 功能特性

- 🗺️ 支持多种坐标系转换：WGS84、GCJ02、BD09、BD09MC、EPSG3857
- 🚀 高性能：单次转换约 100-150ns（Apple M1 Pro）
- 📦 零依赖：仅使用 Go 标准库
- 🎯 高精度：经纬度转换精度约 1 米，投影坐标精度约 1 米
- 🔄 支持多种输入格式：Position、GeoJSON、JSON 字符串
- ✅ 全面测试：基于真实城市数据验证，覆盖全量互转组合

## 支持的坐标系

| 坐标系 | 说明 | 别名 |
|--------|------|------|
| WGS84 | 世界大地坐标系，GPS 原始坐标 | WGS1984, EPSG4326 |
| GCJ02 | 国测局坐标系，中国标准，火星坐标 | AMap |
| BD09 | 百度坐标系，百度地图使用 | BD09LL, Baidu, BMap |
| BD09MC | 百度墨卡托投影坐标系 | BD09Meter |
| EPSG3857 | Web 墨卡托投影坐标系，Google Maps 等使用 | EPSG900913, EPSG102100, WebMercator, WM |

## 安装

```bash
go get github.com/bytebotgo/gcoord-go
```

## 快速开始

### 基础坐标转换

```go
package main

import (
    "fmt"
    "github.com/bytebotgo/gcoord-go"
)

func main() {
    // WGS84 转 GCJ02
    wgs84 := gcoord.Position{116.397, 39.908} // 北京天安门
    gcj02, err := gcoord.Transform(wgs84, gcoord.WGS84, gcoord.GCJ02)
    if err != nil {
        panic(err)
    }
    fmt.Printf("WGS84: %v\n", wgs84)
    fmt.Printf("GCJ02: %v\n", gcj02)
    
    // GCJ02 转 BD09
    bd09, err := gcoord.Transform(gcj02, gcoord.GCJ02, gcoord.BD09)
    if err != nil {
        panic(err)
    }
    fmt.Printf("BD09:  %v\n", bd09)
}
```

### GeoJSON 转换

```go
// 转换 GeoJSON 对象
feature := map[string]any{
    "type": "Feature",
    "geometry": map[string]any{
        "type": "Point",
        "coordinates": []any{116.397, 39.908},
    },
    "properties": map[string]any{
        "name": "天安门",
    },
}

converted, err := gcoord.Transform(feature, gcoord.WGS84, gcoord.GCJ02)
if err != nil {
    panic(err)
}
fmt.Printf("转换后: %+v\n", converted)
```

### JSON 字符串转换

```go
// 转换 JSON 字符串
jsonStr := `{
    "type": "FeatureCollection",
    "features": [{
        "type": "Feature",
        "geometry": {
            "type": "Point",
            "coordinates": [116.397, 39.908]
        }
    }]
}`

result, err := gcoord.Transform(jsonStr, gcoord.WGS84, gcoord.BD09)
if err != nil {
    panic(err)
}
fmt.Printf("转换结果: %s\n", result)
```

## API 参考

### 类型定义

```go
// 坐标系类型
type CRSTypes string

// 坐标位置 [经度, 纬度] 或 [x, y]
type Position []float64
```

### 主要函数

```go
// Transform 执行坐标系转换
func Transform[T any](input T, crsFrom, crsTo CRSTypes) (T, error)
```

**参数：**
- `input`: 输入数据，支持以下类型：
  - `Position`: 坐标数组
  - `string`: JSON 字符串
  - `map[string]any`: GeoJSON 对象
  - `[]any`: 坐标数组
- `crsFrom`: 源坐标系
- `crsTo`: 目标坐标系

**返回值：**
- 转换后的数据（类型与输入相同）
- 错误信息

## 精度说明

- **经纬度转换精度**：约 1e-5 度（约 1 米）
- **投影坐标转换精度**：约 1 米
- **WGS84↔GCJ02**：在中国境外无偏移，直接返回原坐标

## 性能基准

在 Apple M1 Pro 上的性能表现：

```
BenchmarkTransform_WGS84_GCJ02-8      	 7980267	       150.9 ns/op	      56 B/op	       3 allocs/op
BenchmarkTransform_GCJ02_BD09-8       	11908320	       102.0 ns/op	      56 B/op	       3 allocs/op
BenchmarkTransform_WGS84_EPSG3857-8   	11908226	       104.5 ns/op	      56 B/op	       3 allocs/op
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行基准测试
go test -bench=. -benchmem ./gcoord

# 运行特定测试
go test -run TestFixtures ./gcoord
```

## 许可证

MIT License

## 参考

- [gcoord-js](https://github.com/hujiulong/gcoord) - JavaScript 版本
- [国测局坐标系](https://zh.wikipedia.org/wiki/GCJ-02) - GCJ02 坐标系说明
- [百度坐标系](https://lbsyun.baidu.com/index.php?title=coordinate) - BD09 坐标系说明
