package gcoord

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type cityRecord struct {
	Attrs  map[string]any       `json:"attrs"`
	Coords map[string][]float64 `json:"coords"`
}

func mustReadFixture(t *testing.T) []cityRecord {
	t.Helper()
	// 使用仓库内相对位置
	p := filepath.Join("..", "gcoord-js", "test", "fixtures", "china-cities.json")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read fixture error: %v", err)
	}
	var recs []cityRecord
	if err := json.Unmarshal(b, &recs); err != nil {
		t.Fatalf("unmarshal fixture error: %v", err)
	}
	return recs
}

func TestFixtures_BasicPairs(t *testing.T) {
	recs := mustReadFixture(t)
	for i := 0; i < len(recs); i++ {
		c := recs[i]
		wgs := Position(c.Coords["WGS84"])
		gcj := Position(c.Coords["GCJ02"])
		bd9 := Position(c.Coords["BD09"])
		m3857 := Position(c.Coords["EPSG3857"])
		bdmc := Position(c.Coords["BD09MC"])

		// WGS->GCJ 与基准比较
		g1, err := Transform(wgs, WGS84, GCJ02)
		if err != nil {
			t.Fatalf("wgs->gcj err: %v", err)
		}
		if !approxPos(g1, gcj, 1e-5) {
			t.Fatalf("wgs->gcj mismatch: got %v want %v (i=%d)", g1, gcj, i)
		}

		// GCJ->BD09 与基准比较
		b1, err := Transform(gcj, GCJ02, BD09)
		if err != nil {
			t.Fatalf("gcj->bd err: %v", err)
		}
		if !approxPos(b1, bd9, 1e-5) {
			t.Fatalf("gcj->bd mismatch: got %v want %v (i=%d)", b1, bd9, i)
		}

		// WGS<->EPSG3857
		m1, err := Transform(wgs, WGS84, EPSG3857)
		if err != nil {
			t.Fatalf("wgs->3857 err: %v", err)
		}
		if !approxPos(m1, m3857, TestPrecisionProjection) { // 米级精度，1m 容差
			t.Fatalf("wgs->3857 mismatch: got %v want %v (i=%d)", m1, m3857, i)
		}
		w1, err := Transform(m1, EPSG3857, WGS84)
		if err != nil {
			t.Fatalf("3857->wgs err: %v", err)
		}
		if !approxPos(w1, wgs, TestPrecisionLonLat) {
			t.Fatalf("3857->wgs mismatch: got %v want %v (i=%d)", w1, wgs, i)
		}

		// BD09 <-> BD09MC
		mc, err := Transform(bd9, BD09, BD09MC)
		if err != nil {
			t.Fatalf("bd->bdmc err: %v", err)
		}
		if !approxPos(mc, bdmc, TestPrecisionProjection) { // 米级精度，1m 容差
			t.Fatalf("bd->bdmc mismatch: got %v want %v (i=%d)", mc, bdmc, i)
		}
		bd, err := Transform(mc, BD09MC, BD09)
		if err != nil {
			t.Fatalf("bdmc->bd err: %v", err)
		}
		if !approxPos(bd, bd9, TestPrecisionLonLat) {
			t.Fatalf("bdmc->bd mismatch: got %v want %v (i=%d)", bd, bd9, i)
		}
		// 额外组合：WGS84->BD09 与基准
		b2, err := Transform(wgs, WGS84, BD09)
		if err != nil {
			t.Fatalf("wgs->bd09 err: %v", err)
		}
		if !approxPos(b2, bd9, 1e-5) {
			t.Fatalf("wgs->bd09 mismatch: got %v want %v (i=%d)", b2, bd9, i)
		}

		// 额外组合：GCJ02 <-> EPSG3857 往返
		m2, err := Transform(gcj, GCJ02, EPSG3857)
		if err != nil {
			t.Fatalf("gcj->3857 err: %v", err)
		}
		g2, err := Transform(m2, EPSG3857, GCJ02)
		if err != nil {
			t.Fatalf("3857->gcj err: %v", err)
		}
		if !approxPos(g2, gcj, 1e-5) {
			t.Fatalf("gcj<->3857 roundtrip mismatch: got %v want %v (i=%d)", g2, gcj, i)
		}
	}
}
