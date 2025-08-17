package main

/***
- 假设有一个名为 students 的表，
包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
  - 要求 ：
    - 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
    - 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
    - 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
    - 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
**/

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	//开启数据库连接
	dsn := "root:asdf@tcp(localhost:3306)/go_study?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		fmt.Println("连接数据库异常")
	}

	//创建表结构
	db.AutoMigrate(&Student{})

	//插入学生数据
	db.Create(&Student{
		Name:  "张三",
		Age:   20,
		Grade: "三年级",
	})

	//插入多个学生信息
	stu := []Student{
		{Name: "小明", Age: 15, Grade: "二年级"},
		{Name: "小红", Age: 21, Grade: "四年级"},
		{Name: "小里", Age: 22, Grade: "五年级"},
		{Name: "张三", Age: 12, Grade: "三年级"},
		{Name: "小刘", Age: 23, Grade: "六年级"},
	}
	db.Create(stu)

	//查询年龄大于18的学生
	var stus []Student
	result := db.Where("Age > ?", 18).Find(&stus)
	if result.Error != nil {
		fmt.Println("查询大于18岁学生信息异常：", result.Error.Error())
		return
	}
	for _, stden := range stus {
		fmt.Println("查询大于18岁学生信息：", stden)
	}

	//将 students 表中姓名为 "张三" 的学生年级更新为 "四年级",修改数据需要  model(),不然不起作用
	upResult := db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
	fmt.Println("更新条数：", upResult.RowsAffected)

	//删除 students 表中年龄小于 15 岁的学生记录
	delstu := db.Where("age < ?", 15).Delete(&Student{})
	fmt.Println("删除条数：", delstu.RowsAffected)

	fmt.Println("===执行结束====")
}

type Student struct {
	ID    int
	Name  string
	Age   int
	Grade string
}
