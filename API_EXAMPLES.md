# GCoord Go API 使用示例

## 基础用法

### 1. 简单坐标转换

```go
package main

import (
    "fmt"
    "github.com/bytebotgo/gcoord-go/gcoord"
)

func main() {
    // WGS84 转 GCJ02
    coord := gcoord.Position{116.397, 39.908} // 北京天安门
    result, err := gcoord.Transform(coord, gcoord.WGS84, gcoord.GCJ02)
    if err != nil {
        panic(err)
    }
    fmt.Printf("转换结果: %.6f, %.6f\n", result[0], result[1])
}
```

### 2. 使用转换器接口

```go
package main

import (
    "fmt"
    "github.com/bytebotgo/gcoord-go/gcoord"
)

func main() {
    // 创建转换器
    converter, err := gcoord.NewConverter(gcoord.WGS84, gcoord.BD09)
    if err != nil {
        panic(err)
    }
    
    // 执行转换
    coord := gcoord.Position{116.397, 39.908}
    result, err := converter.Convert(coord)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("精度: %f\n", converter.GetPrecision())
    fmt.Printf("结果: %.6f, %.6f\n", result[0], result[1])
}
```

### 3. GeoJSON 转换

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/bytebotgo/gcoord-go/gcoord"
)

func main() {
    // GeoJSON Feature
    feature := map[string]interface{}{
        "type": "Feature",
        "geometry": map[string]interface{}{
            "type": "Point",
            "coordinates": []float64{116.397, 39.908},
        },
        "properties": map[string]interface{}{
            "name": "天安门",
        },
    }
    
    // 转换整个 Feature
    result, err := gcoord.Transform(feature, gcoord.WGS84, gcoord.GCJ02)
    if err != nil {
        panic(err)
    }
    
    // 输出结果
    jsonData, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(jsonData))
}
```

### 4. 批量转换

```go
package main

import (
    "fmt"
    "github.com/bytebotgo/gcoord-go/gcoord"
)

func main() {
    coords := []gcoord.Position{
        {116.397, 39.908}, // 北京
        {121.473, 31.230}, // 上海
        {113.264, 23.129}, // 广州
    }
    
    for i, coord := range coords {
        result, err := gcoord.Transform(coord, gcoord.WGS84, gcoord.GCJ02)
        if err != nil {
            fmt.Printf("转换失败 %d: %v\n", i, err)
            continue
        }
        fmt.Printf("城市 %d: %.6f, %.6f\n", i+1, result[0], result[1])
    }
}
```

### 5. 错误处理

```go
package main

import (
    "fmt"
    "github.com/bytebotgo/gcoord-go/gcoord"
)

func main() {
    coord := gcoord.Position{116.397, 39.908}
    
    result, err := gcoord.Transform(coord, gcoord.WGS84, gcoord.CRSTypes("INVALID"))
    if err != nil {
        if gcoord.IsTransformError(err) {
            errorType := gcoord.GetErrorType(err)
            switch errorType {
            case gcoord.ErrInvalidCRS:
                fmt.Println("无效的坐标系")
            case gcoord.ErrInvalidInput:
                fmt.Println("无效的输入")
            default:
                fmt.Printf("转换错误: %v\n", err)
            }
        } else {
            fmt.Printf("其他错误: %v\n", err)
        }
        return
    }
    
    fmt.Printf("转换成功: %.6f, %.6f\n", result[0], result[1])
}
```

## 高级用法

### 1. 自定义转换器

```go
package main

import (
    "fmt"
    "github.com/bytebotgo/gcoord-go/gcoord"
)

// 自定义转换器
type CustomConverter struct {
    source gcoord.CRSTypes
    target gcoord.CRSTypes
}

func (c *CustomConverter) Convert(coord gcoord.Position) (gcoord.Position, error) {
    // 自定义转换逻辑
    return gcoord.Position{coord[0] + 0.001, coord[1] + 0.001}, nil
}

func (c *CustomConverter) GetSourceCRS() gcoord.CRSTypes {
    return c.source
}

func (c *CustomConverter) GetTargetCRS() gcoord.CRSTypes {
    return c.target
}

func (c *CustomConverter) GetPrecision() float64 {
    return 0.001
}

func main() {
    converter := &CustomConverter{
        source: gcoord.WGS84,
        target: gcoord.CRSTypes("CUSTOM"),
    }
    
    coord := gcoord.Position{116.397, 39.908}
    result, err := converter.Convert(coord)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("自定义转换: %.6f, %.6f\n", result[0], result[1])
}
```

### 2. 性能测试

```go
package main

import (
    "fmt"
    "time"
    "github.com/bytebotgo/gcoord-go/gcoord"
)

func main() {
    coord := gcoord.Position{116.397, 39.908}
    iterations := 100000
    
    // 测试转换性能
    start := time.Now()
    for i := 0; i < iterations; i++ {
        _, _ = gcoord.Transform(coord, gcoord.WGS84, gcoord.GCJ02)
    }
    duration := time.Since(start)
    
    fmt.Printf("转换 %d 次耗时: %v\n", iterations, duration)
    fmt.Printf("平均每次: %v\n", duration/time.Duration(iterations))
}
```

## 最佳实践

1. **错误处理**: 始终检查转换函数的返回值
2. **精度考虑**: 根据应用场景选择合适的精度
3. **性能优化**: 对于批量转换，考虑使用转换器接口
4. **坐标验证**: 转换前验证输入坐标的有效性
5. **内存管理**: 对于大量数据，考虑流式处理

## 性能基准

基于 Apple M1 Pro 的基准测试结果：

```
BenchmarkTransform_WGS84_GCJ02-8      	 4024905	       297.8 ns/op
BenchmarkTransform_GCJ02_BD09-8       	 4938212	       242.3 ns/op
BenchmarkTransform_WGS84_EPSG3857-8   	 4894712	       243.0 ns/op
```

转换器缓存机制显著提升了性能，特别是在重复转换相同坐标系对时。
