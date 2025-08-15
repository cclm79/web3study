package ptr

import "fmt"

func addTen(ptr *int) {
	*ptr += 10 // 解引用并增加10
}

func TestFirst() {
	num := 5
	fmt.Println("修改前:", num)
	addTen(&num) // 传递num的指针
	fmt.Println("修改后:", num)
}
