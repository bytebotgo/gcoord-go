package gcoord

import "math"

// 使用 constants.go 中定义的常量

// BD09ToGCJ02 百度经纬度转火星坐标
func BD09ToGCJ02(coord Position) Position {
	lon, lat := coord[0], coord[1]
	x := lon - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*BaiduFactor)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*BaiduFactor)
	newLon := z * math.Cos(theta)
	newLat := z * math.Sin(theta)
	return Position{newLon, newLat}
}

// GCJ02ToBD09 火星坐标转百度经纬度
func GCJ02ToBD09(coord Position) Position {
	lon, lat := coord[0], coord[1]
	x := lon
	y := lat
	z := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*BaiduFactor)
	theta := math.Atan2(y, x) + 0.000003*math.Cos(x*BaiduFactor)
	newLon := z*math.Cos(theta) + 0.0065
	newLat := z*math.Sin(theta) + 0.006
	return Position{newLon, newLat}
}
