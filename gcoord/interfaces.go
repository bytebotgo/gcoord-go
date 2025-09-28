package gcoord

// CoordinateConverter 坐标转换器接口
type CoordinateConverter interface {
	// Convert 执行坐标转换
	Convert(coord Position) (Position, error)

	// GetSourceCRS 获取源坐标系
	GetSourceCRS() CRSTypes

	// GetTargetCRS 获取目标坐标系
	GetTargetCRS() CRSTypes

	// GetPrecision 获取转换精度
	GetPrecision() float64
}

// GeoJSONProcessor GeoJSON处理器接口
type GeoJSONProcessor interface {
	// Process 处理GeoJSON对象
	Process(obj interface{}) (interface{}, error)

	// CanProcess 检查是否可以处理该对象
	CanProcess(obj interface{}) bool
}

// TransformValidator 转换验证器接口
type TransformValidator interface {
	// ValidateInput 验证输入
	ValidateInput(input interface{}) error

	// ValidateCRS 验证坐标系
	ValidateCRS(crs CRSTypes) error

	// ValidateResult 验证结果
	ValidateResult(result interface{}) error
}

// 基础转换器实现
type baseConverter struct {
	sourceCRS CRSTypes
	targetCRS CRSTypes
	converter Converter
	precision float64
}

func (c *baseConverter) Convert(coord Position) (Position, error) {
	if err := validatePosition(coord); err != nil {
		return nil, err
	}
	return c.converter(coord), nil
}

func (c *baseConverter) GetSourceCRS() CRSTypes {
	return c.sourceCRS
}

func (c *baseConverter) GetTargetCRS() CRSTypes {
	return c.targetCRS
}

func (c *baseConverter) GetPrecision() float64 {
	return c.precision
}

// NewConverter 创建新的转换器
func NewConverter(from, to CRSTypes) (CoordinateConverter, error) {
	converter := getConverter(from, to)
	if converter == nil {
		return nil, ErrUnsupportedCRS(to)
	}

	precision := LonLatPrecision
	if to == EPSG3857 || to == BD09MC {
		precision = ProjectionPrecision
	}

	return &baseConverter{
		sourceCRS: from,
		targetCRS: to,
		converter: converter,
		precision: precision,
	}, nil
}
