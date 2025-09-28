package gcoord

// ensureNumberSlice 确保 Position 至少包含两个元素
func ensureNumberSlice(p Position) Position {
	if len(p) >= 2 {
		return Position{p[0], p[1]}
	}
	return p
}

// validatePosition 验证 Position 是否有效
func validatePosition(p Position) error {
	if len(p) < 2 {
		return ErrInvalidPosition
	}
	return nil
}

// validateCRS 验证坐标系是否有效
func validateCRS(crs CRSTypes) error {
	if crs == "" {
		return ErrEmptyCRS
	}

	validCRS := map[CRSTypes]bool{
		WGS84:    true,
		GCJ02:    true,
		BD09:     true,
		BD09MC:   true,
		EPSG3857: true,
	}

	if !validCRS[crs] {
		return ErrUnsupportedCRS(crs)
	}

	return nil
}
