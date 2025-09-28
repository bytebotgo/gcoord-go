# gcoord-go 命令行工具

一个功能强大的地理坐标转换命令行工具，使用 `github.com/spf13/cobra` 构建，支持多种坐标系之间的转换。

## 🚀 快速开始

### 构建

```bash
# 使用 Makefile
make build

# 或直接使用 go build
go build -o gcoord ./cmd/gcoord
```

### 安装

```bash
# 安装到系统
make install

# 或直接使用 go install
go install ./cmd/gcoord
```

## 📖 使用方法

### 基本用法

```bash
# 显示帮助信息
gcoord --help

# 显示版本信息
gcoord --version
```

### 坐标转换

#### 转换单个坐标点

```bash
# 从 WGS84 转换到 GCJ02
gcoord convert --from WGS84 --to GCJ02 --lon 116.397 --lat 39.908

# 简洁输出 (仅显示结果)
gcoord convert --from WGS84 --to GCJ02 --lon 116.397 --lat 39.908
# 输出: 116.403243,39.909403

# 详细输出
gcoord convert --from WGS84 --to GCJ02 --lon 116.397 --lat 39.908 --verbose
```

#### 转换 JSON 格式坐标

```bash
# 转换 GeoJSON Point
gcoord convert --from WGS84 --to BD09 --json '{"type":"Point","coordinates":[116.397,39.908]}'

# 转换 GeoJSON Feature
gcoord convert --from GCJ02 --to EPSG3857 --json '{"type":"Feature","geometry":{"type":"Point","coordinates":[116.397,39.908]}}'
```

### 查看支持的坐标系

```bash
# 显示所有支持的坐标系
gcoord list
```

## 🌍 支持的坐标系

| 坐标系 | 描述 | 别名 | 精度 |
|--------|------|------|------|
| WGS84 | 世界大地坐标系，GPS原始坐标 | WGS1984, EPSG4326 | 约 1e-5 度 (约 1 米) |
| GCJ02 | 国测局坐标系，中国标准坐标系 | AMap | 约 1e-5 度 (约 1 米) |
| BD09 | 百度坐标系 | BD09LL, Baidu, BMap | 约 1e-5 度 (约 1 米) |
| BD09MC | 百度墨卡托投影坐标系 | BD09Meter | 约 1 米 |
| EPSG3857 | Web墨卡托投影坐标系 | EPSG900913, EPSG102100, WebMercator, WM | 约 1 米 |

## 🎯 功能特性

- ✅ **多种输入格式**: 支持经纬度坐标和 JSON 格式
- ✅ **彩色输出**: 使用 `github.com/fatih/color` 提供美观的彩色输出
- ✅ **专业 CLI**: 基于 `github.com/spf13/cobra` 构建的专业命令行界面
- ✅ **详细模式**: 支持简洁和详细两种输出模式
- ✅ **错误处理**: 完善的错误处理和用户友好的提示信息
- ✅ **自动补全**: 支持 shell 自动补全
- ✅ **链式转换**: 自动选择最优转换路径

## 📋 命令参考

### 主命令

```bash
gcoord [command] [flags]
```

### 子命令

#### convert - 坐标转换

```bash
gcoord convert [flags]

Flags:
  -f, --from string   源坐标系 (必需)
  -t, --to string     目标坐标系 (必需)
      --lon float     经度
      --lat float     纬度
  -j, --json string   JSON格式的坐标输入
  -v, --verbose       显示详细信息
  -h, --help          help for convert
```

#### list - 显示支持的坐标系

```bash
gcoord list [flags]

Flags:
  -h, --help   help for list
```

## 🔧 开发

### 项目结构

```
cmd/gcoord/
├── main.go          # 主程序入口
└── ...

gcoord/
├── transform.go     # 坐标转换核心逻辑
├── crs_types.go     # 坐标系类型定义
└── ...
```

### 构建工具

```bash
# 查看所有可用命令
make help

# 构建
make build

# 安装
make install

# 清理
make clean

# 运行测试
make test

# 格式化代码
make fmt

# 检查代码
make vet
```

## 📝 示例

### 示例 1: GPS 坐标转换

```bash
# GPS 原始坐标 (WGS84) 转换为中国标准坐标 (GCJ02)
gcoord convert --from WGS84 --to GCJ02 --lon 116.397 --lat 39.908 --verbose
```

### 示例 2: 批量转换 JSON 数据

```bash
# 转换 GeoJSON FeatureCollection
gcoord convert --from WGS84 --to BD09 --json '{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [116.397, 39.908]
      }
    }
  ]
}'
```

### 示例 3: 投影坐标转换

```bash
# 从经纬度转换为 Web 墨卡托投影
gcoord convert --from WGS84 --to EPSG3857 --lon 116.397 --lat 39.908 --verbose
```

## 🐛 故障排除

### 常见问题

1. **参数解析错误**: 确保使用正确的参数格式，如 `--from` 而不是 `-from`
2. **坐标系无效**: 使用 `gcoord list` 查看支持的坐标系
3. **JSON 格式错误**: 确保 JSON 格式正确，可以使用在线 JSON 验证工具

### 获取帮助

```bash
# 查看主命令帮助
gcoord --help

# 查看子命令帮助
gcoord convert --help
gcoord list --help
```

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。
