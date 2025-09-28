package gcoord

import "math"

// 使用 constants.go 中定义的常量

// EPSG3857ToWGS84 WebMercator -> WGS84
func EPSG3857ToWGS84(xy Position) Position {
	return Position{
		(xy[0] * RadToDeg) / WGS84A,
		(math.Pi*0.5 - 2.0*math.Atan(math.Exp(-xy[1]/WGS84A))) * RadToDeg,
	}
}

// WGS84ToEPSG3857 WGS84 -> WebMercator
func WGS84ToEPSG3857(lonLat Position) Position {
	adjusted := lonLat[0]
	if math.Abs(lonLat[0]) > 180 {
		if lonLat[0] < 0 {
			adjusted = lonLat[0] + 360
		} else {
			adjusted = lonLat[0] - 360
		}
	}
	xy := Position{
		WGS84A * adjusted * DegToRad,
		WGS84A * math.Log(math.Tan(math.Pi*0.25+0.5*lonLat[1]*DegToRad)),
	}
	if xy[0] > MaxExtent {
		xy[0] = MaxExtent
	}
	if xy[0] < -MaxExtent {
		xy[0] = -MaxExtent
	}
	if xy[1] > MaxExtent {
		xy[1] = MaxExtent
	}
	if xy[1] < -MaxExtent {
		xy[1] = -MaxExtent
	}
	return xy
}
