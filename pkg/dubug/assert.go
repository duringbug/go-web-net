package dubug

import (
	"reflect"
	"testing"
)

// Equal 比较两个值是否相等（支持切片比较）
func Equal(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == b
	}

	switch v := a.(type) {
	case []int:
		if bSlice, ok := b.([]int); ok {
			if len(v) != len(bSlice) {
				return false
			}
			for i := range v {
				if v[i] != bSlice[i] {
					return false
				}
			}
			return true
		}
	default:
		return reflect.DeepEqual(a, b) // 对于其他类型使用 reflect.DeepEqual
	}

	return false
}

// NoError 断言错误值是否为 nil
func NoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
