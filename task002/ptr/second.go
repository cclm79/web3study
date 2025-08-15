package ptr

import "fmt"

func doubleSlice(slicePtr *[]int) {
	// 解引用获取原始切片
	s := *slicePtr
	// 遍历并修改每个元素
	for i := range s {
		s[i] *= 2
	}
}

func TestSecond() {
	nums := []int{1, 2, 3, 4}
	fmt.Println("修改前:", nums)
	doubleSlice(&nums) // 传递切片的指针
	fmt.Println("修改后:", nums)
}
