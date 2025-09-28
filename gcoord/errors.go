package gcoord

import (
	"fmt"
)

// ErrorType 错误类型
type ErrorType int

const (
	ErrInvalidCRS ErrorType = iota
	ErrInvalidInput
	ErrTransformFailed
	ErrUnsupportedFormat
)

// TransformError 转换错误
type TransformError struct {
	Type    ErrorType
	Message string
	Details map[string]interface{}
}

func (e *TransformError) Error() string {
	return e.Message
}

// 预定义错误
var (
	ErrEmptyCRS = &TransformError{
		Type:    ErrInvalidCRS,
		Message: "crsFrom 和 crsTo 不能为空",
	}

	ErrInvalidPosition = &TransformError{
		Type:    ErrInvalidInput,
		Message: "Position 长度必须>=2",
	}
)

// ErrUnsupportedCRS 创建不支持的坐标系错误
func ErrUnsupportedCRS(crs CRSTypes) *TransformError {
	return &TransformError{
		Type:    ErrInvalidCRS,
		Message: fmt.Sprintf("不支持的坐标系: %s", crs),
		Details: map[string]interface{}{
			"crs": crs,
		},
	}
}

// ErrJSONParseFailed 创建JSON解析失败错误
func ErrJSONParseFailed(err error) *TransformError {
	return &TransformError{
		Type:    ErrInvalidInput,
		Message: fmt.Sprintf("JSON解析失败: %v", err),
		Details: map[string]interface{}{
			"original_error": err,
		},
	}
}

// IsTransformError 检查是否为转换错误
func IsTransformError(err error) bool {
	_, ok := err.(*TransformError)
	return ok
}

// GetErrorType 获取错误类型
func GetErrorType(err error) ErrorType {
	if te, ok := err.(*TransformError); ok {
		return te.Type
	}
	return -1
}
