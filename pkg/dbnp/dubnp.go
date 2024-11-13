package np

import (
	"errors"
	"fmt"
)

type Array struct {
	Data  []float64 // 存储数据的扁平化数组
	Shape []int     // 数组的形状（维度）
}

// 创建新的矩阵
func NewArray(data []float64, shape []int) (*Array, error) {
	totalSize := 1
	for _, s := range shape {
		totalSize *= s
	}
	if totalSize != len(data) {
		return nil, errors.New("数据大小与形状不匹配")
	}
	return &Array{Data: data, Shape: shape}, nil
}

// 打印矩阵
func (a *Array) Print() {
	fmt.Println("Shape:", a.Shape)
	fmt.Println("Data:", a.Data)
}
