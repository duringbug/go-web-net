package dubnp

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

type Array struct {
	Data  []float64 // 存储数据的扁平化数组
	Shape []int     // 数组的形状（维度）
}

// 获取系统的内存页面大小，单位为字节
func getPageSize() int {
	return int(unsafe.Sizeof(uintptr(0))) * 8 // 假设每个指针的大小是 8 字节
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

// 矩阵加法（并行加速版）
func (a *Array) Add(b *Array) (*Array, error) {
	// 检查形状是否相同
	if len(a.Shape) != len(b.Shape) {
		return nil, errors.New("矩阵的维度不相同")
	}
	for i := range a.Shape {
		if a.Shape[i] != b.Shape[i] {
			return nil, errors.New("矩阵的形状不匹配")
		}
	}

	// 创建结果数组
	resultData := make([]float64, len(a.Data))

	// 设置并行的 goroutine 数量
	numWorkers := runtime.NumCPU()
	chunkSize := (len(a.Data) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// 并行处理每个块
	for worker := 0; worker < numWorkers; worker++ {
		go func(worker int) {
			defer wg.Done()
			start := worker * chunkSize
			end := start + chunkSize
			if end > len(a.Data) {
				end = len(a.Data)
			}

			// 执行加法运算
			for i := start; i < end; i++ {
				resultData[i] = a.Data[i] + b.Data[i]
			}
		}(worker)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 返回新的矩阵
	return &Array{Data: resultData, Shape: a.Shape}, nil
}

// 打印矩阵数据，按维度格式化输出，decimalPlaces控制小数位数
func (a *Array) PrintMatrix(decimalPlaces int) {
	// 计算每一行的元素个数
	dimensions := len(a.Shape)
	elementsPerRow := a.Shape[dimensions-1]

	// 格式化字符串，用来控制小数点后的位数
	formatString := fmt.Sprintf("%%.%df ", decimalPlaces)

	// 对数据进行格式化输出
	for i := 0; i < len(a.Data); i++ {
		// 每行末尾换行
		if i%elementsPerRow == 0 && i != 0 {
			fmt.Println()
		}
		// 打印数据，并用空格分隔
		fmt.Printf(formatString, a.Data[i])
	}
	fmt.Println()
}

// 矩阵乘法（优化版，内存访问优化 + 并行化 + 分块优化）
func (a *Array) Multiply(b *Array) (*Array, error) {
	// 检查矩阵维度是否符合乘法要求
	if len(a.Shape) != 2 || len(b.Shape) != 2 {
		return nil, errors.New("仅支持二维矩阵乘法")
	}
	if a.Shape[1] != b.Shape[0] {
		return nil, errors.New("矩阵的维度不匹配，无法进行乘法运算")
	}

	// 创建结果矩阵
	resultData := make([]float64, a.Shape[0]*b.Shape[1])

	// 转置矩阵 b，优化列访问
	bTransposed, err := b.Transpose()
	if err != nil {
		return nil, err
	}

	// 获取系统的内存页面大小，用于确定合理的块大小
	pageSize := getPageSize()
	blockSize := pageSize // 设置块大小为内存页面大小

	// 设置并行的 goroutine 数量
	numWorkers := runtime.NumCPU()

	// 计算每个块的维度
	blockRows := blockSize / 8 // 每块的行数 (假设每个浮点数占 8 字节)
	blockCols := blockSize / 8 // 每块的列数 (假设每个浮点数占 8 字节)

	// 确保块的维度不超过矩阵的实际维度
	if blockRows > a.Shape[0] {
		blockRows = a.Shape[0]
	}
	if blockCols > b.Shape[1] {
		blockCols = b.Shape[1]
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// 并行处理每个块
	for worker := 0; worker < numWorkers; worker++ {
		go func(worker int) {
			defer wg.Done()

			// 计算当前 goroutine 需要处理的块的起始和结束索引
			for i := worker * blockRows; i < a.Shape[0]; i += blockRows {
				for j := worker * blockCols; j < b.Shape[1]; j += blockCols {
					// 对每个小块进行矩阵计算
					for k := 0; k < a.Shape[1]; k++ {
						// 执行实际的乘法和加法操作
						resultData[i*b.Shape[1]+j] += a.Data[i*a.Shape[1]+k] * bTransposed.Data[j*b.Shape[0]+k]
					}
				}
			}
		}(worker)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 返回新的矩阵
	return &Array{Data: resultData, Shape: []int{a.Shape[0], b.Shape[1]}}, nil
}

// 转置矩阵（并行加速版，简单并行化）
func (a *Array) Transpose() (*Array, error) {
	// 检查矩阵是否为二维
	if len(a.Shape) != 2 {
		return nil, fmt.Errorf("仅支持二维矩阵转置")
	}

	// 创建结果矩阵，大小为 a.Shape[1] x a.Shape[0]
	resultData := make([]float64, a.Shape[1]*a.Shape[0])

	// 设置并行的 goroutine 数量
	numWorkers := runtime.NumCPU()

	// 设置每个 goroutine 需要处理的数据范围
	chunkSize := (a.Shape[0]*a.Shape[1] + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// 并行处理每个块
	for worker := 0; worker < numWorkers; worker++ {
		go func(worker int) {
			defer wg.Done()

			// 计算当前 goroutine 需要处理的起始和结束索引
			start := worker * chunkSize
			end := start + chunkSize
			if end > a.Shape[0]*a.Shape[1] {
				end = a.Shape[0] * a.Shape[1]
			}

			// 执行矩阵转置操作
			for i := start; i < end; i++ {
				// 计算行和列的索引
				row := i / a.Shape[1]
				col := i % a.Shape[1]
				// 转置过程：将 (row, col) 元素转置到 (col, row)
				resultData[col*a.Shape[0]+row] = a.Data[row*a.Shape[1]+col]
			}
		}(worker)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 返回新的转置矩阵
	return &Array{Data: resultData, Shape: []int{a.Shape[1], a.Shape[0]}}, nil
}
