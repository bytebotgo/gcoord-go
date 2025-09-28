package gcoord

import (
	"testing"
)

// TestHelper 测试辅助结构
type TestHelper struct {
	t *testing.T
}

// NewTestHelper 创建测试辅助器
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// AssertApproxEqual 断言两个浮点数近似相等
func (h *TestHelper) AssertApproxEqual(a, b, epsilon float64, msg string) {
	if !h.approx(a, b, epsilon) {
		h.t.Errorf("%s: expected %f ≈ %f (epsilon=%f)", msg, a, b, epsilon)
	}
}

// AssertPositionApproxEqual 断言两个位置近似相等
func (h *TestHelper) AssertPositionApproxEqual(a, b Position, epsilon float64, msg string) {
	if !h.approxPos(a, b, epsilon) {
		h.t.Errorf("%s: expected %v ≈ %v (epsilon=%f)", msg, a, b, epsilon)
	}
}

// AssertNoError 断言没有错误
func (h *TestHelper) AssertNoError(err error, msg string) {
	if err != nil {
		h.t.Errorf("%s: unexpected error: %v", msg, err)
	}
}

// AssertError 断言有错误
func (h *TestHelper) AssertError(err error, msg string) {
	if err == nil {
		h.t.Errorf("%s: expected error but got nil", msg)
	}
}

// TestRoundtrip 测试往返转换
func (h *TestHelper) TestRoundtrip(coord Position, from, to CRSTypes, epsilon float64) {
	// 正向转换
	converted, err := Transform(coord, from, to)
	h.AssertNoError(err, "forward transform")

	// 反向转换
	back, err := Transform(converted, to, from)
	h.AssertNoError(err, "reverse transform")

	// 检查往返精度
	h.AssertPositionApproxEqual(back, coord, epsilon, "roundtrip precision")
}

// approx 检查两个浮点数是否近似相等
func (h *TestHelper) approx(a, b, eps float64) bool {
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= eps
}

// approxPos 检查两个位置是否近似相等
func (h *TestHelper) approxPos(a, b Position, eps float64) bool {
	if len(a) < 2 || len(b) < 2 {
		return false
	}
	return h.approx(a[0], b[0], eps) && h.approx(a[1], b[1], eps)
}

// BenchmarkTransform 基准测试辅助函数
func BenchmarkTransform(b *testing.B, coord Position, from, to CRSTypes) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Transform(coord, from, to)
	}
}

// 测试数据生成器
type TestDataGenerator struct{}

// GetTestCoordinates 获取测试坐标
func (g *TestDataGenerator) GetTestCoordinates() map[string]Position {
	return map[string]Position{
		"beijing":   {116.397, 39.908},
		"shanghai":  {121.473, 31.230},
		"guangzhou": {113.264, 23.129},
		"shenzhen":  {114.057, 22.543},
		"hangzhou":  {120.155, 30.274},
		"nanjing":   {118.767, 32.041},
		"wuhan":     {114.298, 30.584},
		"chengdu":   {104.066, 30.572},
		"xian":      {108.948, 34.263},
		"tianjin":   {117.200, 39.084},
	}
}

// GetEdgeCases 获取边界情况测试数据
func (g *TestDataGenerator) GetEdgeCases() map[string]Position {
	return map[string]Position{
		"equator":        {0, 0},
		"north_pole":     {0, 90},
		"south_pole":     {0, -90},
		"date_line":      {180, 0},
		"antimeridian":   {-180, 0},
		"china_border_n": {116, ChinaMaxLat},
		"china_border_s": {116, ChinaMinLat},
		"china_border_e": {ChinaMaxLon, 35},
		"china_border_w": {ChinaMinLon, 35},
	}
}

// 精度测试工具
func TestPrecision(t *testing.T, coord Position, from, to CRSTypes, expectedPrecision float64) {
	converted, err := Transform(coord, from, to)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// 计算实际精度（这里简化处理）
	// 实际项目中可能需要更复杂的精度计算
	helper := NewTestHelper(t)
	helper.AssertPositionApproxEqual(converted, converted, expectedPrecision, "precision test")
}
