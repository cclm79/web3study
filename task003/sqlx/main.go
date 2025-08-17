package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/** Sqlx入门
1.1 题目1：使用SQL扩展库进行查询
  - 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
    - 要求 ：
      - 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
      - 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
1.2 题目2：实现类型安全映射
  - 假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
    - 要求 ：
      - 定义一个 Book 结构体，包含与 books 表对应的字段。
      - 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

*/

func main() {

	// sqlx 的连接
	db, err := sqlx.Connect("mysql", "root:asdf@tcp(localhost:3306)/go_study?charset=utf8mb4&parseTime=True&loc=Local")

	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 查询技术部门数据
	queryEmployee(db)
	//查询工资最高的员工信息
	queryTopSalary(db)

	//插入数据
	// dsn := "insert into  books(title,author,price)values(?,?,?)"
	// db.Exec(dsn, "秦时明月", "小明", 100.0)
	// db.Exec(dsn, "秦时明月02", "小明01", 2000.0)
	// db.Exec(dsn, "秦时明月03", "小明01", 80.0)
	// 查询数据
	queryBooks(db)

	fmt.Println("===============")

}

func queryBooks(db *sqlx.DB) {
	var boks []Book
	dsn := "select id,title,author,price from books where price > ?"
	err1 := db.Select(&boks, dsn, 50.0)
	if err1 != nil {
		panic("查询数据异常")
	}

	for _, book := range boks {
		fmt.Println(book)
	}
}

func queryEmployee(db *sqlx.DB) {
	// 查询技术部门数据
	var emps []Employee
	dsn := "select id,name,department,salary from employees where department=?"
	err1 := db.Select(&emps, dsn, "技术部")
	if err1 != nil {
		panic("查询数据异常: ")
	}
	fmt.Print("技术部员工：")
	fmt.Println(emps)
}

func queryTopSalary(db *sqlx.DB) {
	var emps Employee
	sql := "select * from employees order by salary desc limit 1"
	err1 := db.Get(&emps, sql)
	if err1 != nil {
		panic("查询数据异常: ")
	}
	fmt.Print("技术部员工：")
	fmt.Println(emps)
}

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
}

type Book struct {
	ID     int
	Title  string
	Author string
	Price  float64
}
