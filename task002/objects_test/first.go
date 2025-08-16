package objects_test

import "fmt"

//
// 考察点 ：接口的定义与实现、面向对象编程风格。
/**
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
*/

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	w int
	h int
}

func (r *Rectangle) Area() int {
	return r.h * r.w
}

func (r *Rectangle) Perimeter() int {
	return (r.h + r.w) * 2
}

type Circle struct {
	r int
}

func (c *Circle) Area() int {
	return c.r * c.r * 3
}

func (c *Circle) Perimeter() int {
	return 2 * c.r * 3
}

func TestFirst() {
	r := Rectangle{h: 2, w: 3}
	fmt.Println(r.Area())
	fmt.Println(r.Perimeter())

	c := Circle{r: 3}
	fmt.Println(c.Area())
	fmt.Println(c.Perimeter())
}
