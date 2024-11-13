package test

import (
	"github.com/duringbug/go-web-net/pkg/dubnp"
	"testing"
)

// 测试 NewArray 函数
func TestNewArray(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		data   []float64
		shape  []int
		hasErr bool // 是否预期返回错误
	}{
		// 合法的矩阵数据和形状
		{[]float64{1, 2, 3, 4}, []int{2, 2}, false},
		// 数据大小与形状不匹配
		{[]float64{1, 2, 3}, []int{2, 2}, true},
	}

	for _, tt := range tests {
		array, err := NewArray(tt.data, tt.shape)
		if (err != nil) != tt.hasErr {
			t.Errorf("NewArray(%v, %v) returned error = %v, expected error = %v", tt.data, tt.shape, err != nil, tt.hasErr)
		}
		if !tt.hasErr && array != nil {
			// 检查返回的数组是否匹配预期
			if len(array.Data) != len(tt.data) || len(array.Shape) != len(tt.shape) {
				t.Errorf("NewArray(%v, %v) created array with wrong data or shape", tt.data, tt.shape)
			}
		}
	}
}

// 测试 Print 方法
func TestArray_Print(t *testing.T) {
	// 创建一个测试用例
	data := []float64{1, 2, 3, 4}
	shape := []int{2, 2}
	array, err := NewArray(data, shape)
	if err != nil {
		t.Fatalf("NewArray failed: %v", err)
	}

	// 捕获 Print 的输出
	array.Print() // 这部分可以通过手动验证输出是否正确
}
