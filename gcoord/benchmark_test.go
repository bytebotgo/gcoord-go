package gcoord

import (
	"testing"
)

func BenchmarkTransform_WGS84_GCJ02(b *testing.B) {
	p := Position{116.397, 39.908}
	for i := 0; i < b.N; i++ {
		_, _ = Transform(p, WGS84, GCJ02)
	}
}

func BenchmarkTransform_GCJ02_BD09(b *testing.B) {
	p := Position{116.404, 39.915}
	for i := 0; i < b.N; i++ {
		_, _ = Transform(p, GCJ02, BD09)
	}
}

func BenchmarkTransform_WGS84_EPSG3857(b *testing.B) {
	p := Position{12.4924, 41.8902}
	for i := 0; i < b.N; i++ {
		_, _ = Transform(p, WGS84, EPSG3857)
	}
}
