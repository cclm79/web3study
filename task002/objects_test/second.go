package objects_test

import "fmt"

// 定义Person结构体
type Person struct {
	Name string
	Age  int
}

// 定义Employee结构体，组合Person
type Employee struct {
	Person     // 匿名嵌入Person实现组合
	EmployeeID string
}

// 为Employee实现PrintInfo方法
func (e Employee) PrintInfo() {
	fmt.Printf("Name: %s\nAge: %d\nEmployee ID: %s\n", e.Name, e.Age, e.EmployeeID)
}

func TestSecond() {
	// 创建Employee实例
	emp := Employee{
		Person: Person{
			Name: "Alice",
			Age:  30,
		},
		EmployeeID: "E12345",
	}

	// 调用PrintInfo方法输出信息
	emp.PrintInfo()
}
