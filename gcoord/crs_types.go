package gcoord

// CRSTypes 表示坐标参考系类型
type CRSTypes string

const (
	// WGS84
	WGS84    CRSTypes = "WGS84"
	WGS1984  CRSTypes = WGS84
	EPSG4326 CRSTypes = WGS84

	// GCJ02
	GCJ02 CRSTypes = "GCJ02"
	AMap  CRSTypes = GCJ02

	// BD09
	BD09   CRSTypes = "BD09"
	BD09LL CRSTypes = BD09
	Baidu  CRSTypes = BD09
	BMap   CRSTypes = BD09

	// BD09MC
	BD09MC    CRSTypes = "BD09MC"
	BD09Meter CRSTypes = BD09MC

	// EPSG3857
	EPSG3857    CRSTypes = "EPSG3857"
	EPSG900913  CRSTypes = EPSG3857
	EPSG102100  CRSTypes = EPSG3857
	WebMercator CRSTypes = EPSG3857
	WM          CRSTypes = EPSG3857
)

// Position 为经纬度或投影坐标 [x, y]，允许长度>=2
type Position []float64
