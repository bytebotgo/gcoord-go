package gcoord

import (
	"encoding/json"
	"math"
	"testing"
)

func approx(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

func approxPos(a, b Position, eps float64) bool {
	if len(a) < 2 || len(b) < 2 {
		return false
	}
	return approx(a[0], b[0], eps) && approx(a[1], b[1], eps)
}

func TestIdentity(t *testing.T) {
	p := Position{116.397, 39.908}
	out, err := Transform(p, WGS84, WGS84)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !approxPos(out, p, 0) {
		t.Fatalf("identity failed: %v != %v", out, p)
	}
}

func TestWGS84GCJ02Roundtrip(t *testing.T) {
	src := Position{116.397, 39.908}
	to, err := Transform(src, WGS84, GCJ02)
	if err != nil {
		t.Fatalf("wgs->gcj error: %v", err)
	}
	back, err := Transform(to, GCJ02, WGS84)
	if err != nil {
		t.Fatalf("gcj->wgs error: %v", err)
	}
	if !approxPos(back, src, TestPrecisionRoundtrip*10) { // 放宽到 1e-5 量级
		t.Fatalf("roundtrip mismatch: got %v want %v", back, src)
	}
}

func TestGCJ02BD09Roundtrip(t *testing.T) {
	src := Position{116.404, 39.915}
	to, err := Transform(src, GCJ02, BD09)
	if err != nil {
		t.Fatalf("gcj->bd09 error: %v", err)
	}
	back, err := Transform(to, BD09, GCJ02)
	if err != nil {
		t.Fatalf("bd09->gcj error: %v", err)
	}
	if !approxPos(back, src, TestPrecisionRoundtrip) { // 约 1e-6 容差
		t.Fatalf("roundtrip mismatch: got %v want %v", back, src)
	}
}

func TestBD09LLBD09MCRoundtrip(t *testing.T) {
	src := Position{116.404, 39.915}
	to, err := Transform(src, BD09, BD09MC)
	if err != nil {
		t.Fatalf("bd09->bd09mc error: %v", err)
	}
	back, err := Transform(to, BD09MC, BD09)
	if err != nil {
		t.Fatalf("bd09mc->bd09 error: %v", err)
	}
	if !approxPos(back, src, TestPrecisionRoundtrip) {
		t.Fatalf("roundtrip mismatch: got %v want %v", back, src)
	}
}

func TestWGS84EPSG3857Roundtrip(t *testing.T) {
	src := Position{12.4924, 41.8902} // 罗马斗兽场附近，经度不跨 180
	to, err := Transform(src, WGS84, EPSG3857)
	if err != nil {
		t.Fatalf("wgs->3857 error: %v", err)
	}
	back, err := Transform(to, EPSG3857, WGS84)
	if err != nil {
		t.Fatalf("3857->wgs error: %v", err)
	}
	if !approxPos(back, src, TestPrecisionRoundtrip) {
		t.Fatalf("roundtrip mismatch: got %v want %v", back, src)
	}
}

func TestTransformJSONStringGeoJSON(t *testing.T) {
	feature := map[string]any{
		"type": "Feature",
		"geometry": map[string]any{
			"type":        "Point",
			"coordinates": []any{116.397, 39.908},
		},
		"properties": map[string]any{"name": "t"},
	}
	b, _ := json.Marshal(feature)
	out, err := Transform(string(b), WGS84, GCJ02)
	if err != nil {
		t.Fatalf("json transform error: %v", err)
	}
	// 只检查结构仍是 JSON 字符串
	if _, ok := any(out).(string); !ok {
		t.Fatalf("expect string output for string input")
	}
}

func TestPreserveExtraDimensions(t *testing.T) {
	// 带第三维与第四维的数据，仅前两维应被转换；额外维度保持数值以符合坐标数组约定
	src := []any{116.397, 39.908, 123.45, 9999.0}
	m := map[string]any{
		"type":        "Point",
		"coordinates": src,
	}
	outAny := transformAny(m, func(p Position) Position { return WGS84ToGCJ02(p) })
	out := outAny.(map[string]any)
	coords := out["coordinates"].([]any)

	if len(coords) != 4 {
		t.Fatalf("expect 4 elements, got %d", len(coords))
	}
	// 前两维应变化（中国境内点），后两维保持原值
	if coords[2] != src[2] || coords[3] != src[3] {
		t.Fatalf("extra dimensions not preserved: got %v want %v", coords[2:], src[2:])
	}
	// 粗检前两维是否已变化
	if coords[0] == src[0] || coords[1] == src[1] {
		t.Fatalf("first two dimensions not transformed: got %v from %v", coords[:2], src[:2])
	}
}

func TestGeometryCollectionTransform(t *testing.T) {
	gc := map[string]any{
		"type": "GeometryCollection",
		"geometries": []any{
			map[string]any{
				"type":        "Point",
				"coordinates": []any{116.397, 39.908},
			},
			map[string]any{
				"type": "LineString",
				"coordinates": []any{
					[]any{116.39, 39.90},
					[]any{116.40, 39.91},
				},
			},
		},
	}

	out := transformAny(gc, func(p Position) Position { return WGS84ToGCJ02(p) }).(map[string]any)
	geoms := out["geometries"].([]any)

	// 检查 Point 已变化
	p := geoms[0].(map[string]any)["coordinates"].([]any)
	if approx(p[0].(float64), 116.397, 0) && approx(p[1].(float64), 39.908, 0) {
		t.Fatalf("Point not transformed: %v", p)
	}
	// 检查 LineString 各点已变化
	line := geoms[1].(map[string]any)["coordinates"].([]any)
	for i := range line {
		pt := line[i].([]any)
		if approx(pt[0].(float64), []float64{116.39, 116.40}[i], 0) || approx(pt[1].(float64), []float64{39.90, 39.91}[i], 0) {
			t.Fatalf("LineString point not transformed: %v", pt)
		}
	}
}
