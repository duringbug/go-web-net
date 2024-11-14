package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/duringbug/go-web-net/pkg/dubnp"
	"github.com/duringbug/go-web-net/pkg/dubug"
)

// 测试 NewArray 函数
func TestNewFloatArray(t *testing.T) {
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
		// 空数据和形状
		{[]float64{}, []int{0, 0}, false},
		// 空数据但不匹配形状
		{[]float64{}, []int{2, 2}, true},
		// 空数据形状为负
		{[]float64{}, []int{-1, -1}, true},
	}

	for _, tt := range tests {
		// 调用 NewArray 函数
		array, err := dubnp.NewArray(tt.data, tt.shape)

		// 检查错误是否符合预期
		if (err != nil) != tt.hasErr {
			t.Errorf("NewArray(%v, %v) returned error = %v, expected error = %v", tt.data, tt.shape, err != nil, tt.hasErr)
		}

		// 如果不应该有错误，则检查返回的 array 是否匹配预期
		if !tt.hasErr {
			if array == nil {
				t.Errorf("NewArray(%v, %v) returned nil array, but expected a valid array", tt.data, tt.shape)
			} else {
				// 检查返回的数组是否匹配预期的数据和形状
				if len(array.Data) != len(tt.data) || len(array.Shape) != len(tt.shape) {
					t.Errorf("NewArray(%v, %v) created array with wrong data or shape. Got data: %v, shape: %v, expected data: %v, shape: %v",
						tt.data, tt.shape, array.Data, array.Shape, tt.data, tt.shape)
				}
			}
		} else {
			// 如果预期有错误，检查 array 是否为 nil
			if array != nil {
				t.Errorf("NewArray(%v, %v) returned a non-nil array, but expected an error", tt.data, tt.shape)
			}
		}
	}
}

// 测试 Print 方法
func TestArrayPrint(t *testing.T) {
	// 创建一个测试用例
	data := []float64{1, 2, 3, 4}
	shape := []int{2, 2}
	array, err := dubnp.NewArray(data, shape)
	if err != nil {
		t.Fatalf("NewArray failed: %v", err)
	}

	// 捕获 Print 的输出
	array.Print() // 这部分可以通过手动验证输出是否正确
	array.PrintMatrix(8)
}

func TestAdd(t *testing.T) {
	// 创建两个矩阵
	a, err := dubnp.NewArray([]float64{1, 2, 3, 4}, []int{2, 2})
	if err != nil {
		t.Fatalf("创建矩阵 a 时出错: %v", err)
	}
	b, err := dubnp.NewArray([]float64{5, 6, 7, 8}, []int{2, 2})
	if err != nil {
		t.Fatalf("创建矩阵 b 时出错: %v", err)
	}

	// 捕获开始时间
	startTime := time.Now()
	// 计算矩阵加法
	result, err := a.Add(b)
	// 计算运行时长
	duration := time.Since(startTime)
	// 打印时长，保留 6 位小数
	fmt.Printf("TestAdd duration: %.9f seconds\n", duration.Seconds())

	if err != nil {
		t.Fatalf("矩阵加法时出错: %v", err)
	}

	// 检查结果的数据是否正确
	expectedData := []float64{6, 8, 10, 12}
	for i, v := range result.Data {
		if v != expectedData[i] {
			t.Errorf("期望数据 %v, 但实际数据为 %v", expectedData, result.Data)
			break
		}
	}
}

// 测试矩阵相加
func TestAddRandomMatrices(t *testing.T) {
	// 使用当前时间戳作为种子来创建一个新的随机数生成器
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// 创建一个 1024x1024 大小的矩阵，元素为随机浮动数
	size := 1024
	dataA := make([]float64, size*size)
	dataB := make([]float64, size*size)

	// 填充数据 A 和 B
	for i := 0; i < size*size; i++ {
		dataA[i] = r.Float64() * 100 // 生成 0 到 100 之间的随机浮动数
		dataB[i] = r.Float64() * 100
	}

	// 创建矩阵 A 和 B
	a, err := dubnp.NewArray(dataA, []int{size, size})
	if err != nil {
		t.Fatalf("创建矩阵 a 时出错: %v", err)
	}
	b, err := dubnp.NewArray(dataB, []int{size, size})
	if err != nil {
		t.Fatalf("创建矩阵 b 时出错: %v", err)
	}

	startTime := time.Now()
	// 计算矩阵加法
	result, err := a.Add(b)
	// 计算运行时长
	duration := time.Since(startTime)
	// 打印时长，保留 6 位小数
	fmt.Printf("TestAddRandomMatrices duration: %.9f seconds\n", duration.Seconds())

	if err != nil {
		t.Fatalf("矩阵加法时出错: %v", err)
	}

	// 检查结果的数据是否正确，示例中只打印了部分结果
	// 打印前几个元素以检查
	// fmt.Println("Matrix A first 5 elements:", dataA[:5])
	// fmt.Println("Matrix B first 5 elements:", dataB[:5])
	// fmt.Println("Result first 5 elements:", result.Data[:5])
	// 断言矩阵加法结果的形状与输入矩阵相同
	if !dubug.Equal(a.Shape, result.Shape) {
		t.Fatalf("矩阵形状不匹配: got %v, want %v", result.Shape, a.Shape)
	}
}

// 测试矩阵相乘
func TestMultiplyRandomMatrices(t *testing.T) {
	// 使用当前时间戳作为种子来创建一个新的随机数生成器
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// 创建一个 1024x1024 大小的矩阵，元素为随机浮动数
	size := 1024
	dataA := make([]float64, size*size)
	dataB := make([]float64, size*size)

	// 填充数据 A 和 B
	for i := 0; i < size*size; i++ {
		dataA[i] = r.Float64() * 100 // 生成 0 到 100 之间的随机浮动数
		dataB[i] = r.Float64() * 100
	}

	// 创建矩阵 A 和 B
	a, err := dubnp.NewArray(dataA, []int{size, size})
	if err != nil {
		t.Fatalf("创建矩阵 a 时出错: %v", err)
	}
	b, err := dubnp.NewArray(dataB, []int{size, size})
	if err != nil {
		t.Fatalf("创建矩阵 b 时出错: %v", err)
	}

	startTime := time.Now()
	// 计算矩阵加法
	result, err := a.Multiply(b)
	// 计算运行时长
	duration := time.Since(startTime)
	// 打印时长，保留 6 位小数
	fmt.Printf("TestMultiplyRandomMatrices duration: %.9f seconds\n", duration.Seconds())

	if err != nil {
		t.Fatalf("矩阵加法时出错: %v", err)
	}

	// 检查结果的数据是否正确，示例中只打印了部分结果
	// 打印前几个元素以检查
	// fmt.Println("Matrix A first 5 elements:", dataA[:5])
	// fmt.Println("Matrix B first 5 elements:", dataB[:5])
	// fmt.Println("Result first 5 elements:", result.Data[:5])
	// 断言矩阵加法结果的形状与输入矩阵相同
	if !dubug.Equal(a.Shape, result.Shape) {
		t.Fatalf("矩阵形状不匹配: got %v, want %v", result.Shape, a.Shape)
	}

}

// 测试矩阵转置
func TestTransposeRandomMatrices(t *testing.T) {
	// 使用当前时间戳作为种子来创建一个新的随机数生成器
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// 创建一个 1024x1024 大小的矩阵，元素为随机浮动数
	size := 1024
	dataA := make([]float64, size*size)

	// 填充数据 A
	for i := 0; i < size*size; i++ {
		dataA[i] = r.Float64() * 100 // 生成 0 到 100 之间的随机浮动数
	}

	// 创建矩阵 A
	a, err := dubnp.NewArray(dataA, []int{size, size})
	if err != nil {
		t.Fatalf("创建矩阵 a 时出错: %v", err)
	}

	startTime := time.Now()
	// 计算矩阵加法
	result, err := a.Transpose()
	// 计算运行时长
	duration := time.Since(startTime)
	// 打印时长，保留 6 位小数
	fmt.Printf("TestTransposeRandomMatrices duration: %.9f seconds\n", duration.Seconds())

	if err != nil {
		t.Fatalf("矩阵加法时出错: %v", err)
	}

	// 检查结果的数据是否正确，示例中只打印了部分结果
	// 打印前几个元素以检查
	// fmt.Println("Matrix A first 5 elements:", dataA[:5])
	// fmt.Println("Result first 5 elements:", result.Data[:5])
	// 断言矩阵加法结果的形状与输入矩阵相同
	if !dubug.Equal(a.Shape, result.Shape) {
		t.Fatalf("矩阵形状不匹配: got %v, want %v", result.Shape, a.Shape)
	}

}
