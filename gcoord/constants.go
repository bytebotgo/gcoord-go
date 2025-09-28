package gcoord

import "math"

// 精度常量
const (
	// 经纬度转换精度
	LonLatPrecision = 1e-5 // 约1米

	// 投影坐标转换精度
	ProjectionPrecision = 1.0 // 1米

	// 迭代收敛精度
	IterationPrecision = 1e-6

	// 中国边界框
	ChinaMinLon = 72.004
	ChinaMaxLon = 137.8347
	ChinaMinLat = 0.8293
	ChinaMaxLat = 55.8271
)

// 地球参数
const (
	// WGS84椭球参数
	WGS84A = 6378137.0
	WGS84F = 1.0 / 298.257223563

	// GCJ02椭球参数
	GCJ02A  = 6378245.0
	GCJ02E2 = 0.006693421622965823

	// 百度坐标系参数
	BaiduFactor = math.Pi * 3000.0 / 180.0

	// 墨卡托投影参数
	MaxExtent = 20037508.342789244
)

// 数学常量
const (
	DegToRad = math.Pi / 180
	RadToDeg = 180 / math.Pi
)

// 转换精度阈值
const (
	// 用于测试的精度阈值
	TestPrecisionLonLat     = 1e-5
	TestPrecisionProjection = 1.0
	TestPrecisionRoundtrip  = 1e-6
)
