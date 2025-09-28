package gcoord

import "math"

// 使用 constants.go 中定义的常量

func isInChinaBbox(lon, lat float64) bool {
	return lon >= ChinaMinLon && lon <= ChinaMaxLon && lat >= ChinaMinLat && lat <= ChinaMaxLat
}

func transformLat(x, y float64) float64 {
	ret := -100 + 2*x + 3*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20*math.Sin(6*x*math.Pi) + 20*math.Sin(2*x*math.Pi)) * 2 / 3
	ret += (20*math.Sin(y*math.Pi) + 40*math.Sin((y/3)*math.Pi)) * 2 / 3
	ret += (160*math.Sin((y/12)*math.Pi) + 320*math.Sin((y*math.Pi)/30)) * 2 / 3
	return ret
}

func transformLon(x, y float64) float64 {
	ret := 300 + x + 2*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20*math.Sin(6*x*math.Pi) + 20*math.Sin(2*x*math.Pi)) * 2 / 3
	ret += (20*math.Sin(x*math.Pi) + 40*math.Sin((x/3)*math.Pi)) * 2 / 3
	ret += (150*math.Sin((x/12)*math.Pi) + 300*math.Sin((x/30)*math.Pi)) * 2 / 3
	return ret
}

func delta(lon, lat float64) (float64, float64) {
	dLon := transformLon(lon-105, lat-35)
	dLat := transformLat(lon-105, lat-35)
	radLat := lat * DegToRad
	magic := math.Sin(radLat)
	magic = 1 - GCJ02E2*magic*magic
	sqrtMagic := math.Sqrt(magic)
	// 注意：这里应使用常数 180 而不是 RadToDeg。原式为：
	// dLon = dLon * 180 / ((a / sqrtMagic) * cos(radLat) * PI)
	// dLat = dLat * 180 / (((a * (1 - ee)) / (magic * sqrtMagic)) * PI)
	dLon = dLon * 180.0 / ((GCJ02A / sqrtMagic) * math.Cos(radLat) * math.Pi)
	dLat = dLat * 180.0 / (((GCJ02A * (1 - GCJ02E2)) / (magic * sqrtMagic)) * math.Pi)
	return dLon, dLat
}

// WGS84ToGCJ02 按 JS 逻辑转换（中国境外不变）
func WGS84ToGCJ02(coord Position) Position {
	lon, lat := coord[0], coord[1]
	if !isInChinaBbox(lon, lat) {
		return Position{lon, lat}
	}
	dLon, dLat := delta(lon, lat)
	return Position{lon + dLon, lat + dLat}
}

// GCJ02ToWGS84 使用迭代反解
func GCJ02ToWGS84(coord Position) Position {
	lon, lat := coord[0], coord[1]
	if !isInChinaBbox(lon, lat) {
		return Position{lon, lat}
	}
	wgsLon, wgsLat := lon, lat
	temp := WGS84ToGCJ02(Position{wgsLon, wgsLat})
	dx := temp[0] - lon
	dy := temp[1] - lat
	for math.Abs(dx) > IterationPrecision || math.Abs(dy) > IterationPrecision {
		wgsLon -= dx
		wgsLat -= dy
		temp = WGS84ToGCJ02(Position{wgsLon, wgsLat})
		dx = temp[0] - lon
		dy = temp[1] - lat
	}
	return Position{wgsLon, wgsLat}
}
