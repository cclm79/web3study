package main

/**
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表
（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
  - 要求 ：
    - 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
      如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	//连接数据库方法
	db, err := ConnenDb()

	if err != nil {
		fmt.Println("连接数据库异常")
		return
	}

	// //创建表结构
	// db.AutoMigrate(&Account{})
	// db.AutoMigrate(&Transaction{})
	// // //插入信息
	// acc := []Account{
	// 	{Name: "小明", Balance: 100.0},
	// 	{Name: "小红", Balance: 90.0},
	// }
	// db.Create(acc)

	// 执行修改
	upadteAmount(1, 2, 30.0, db)

}

// 执行a,b账户调额
func upadteAmount(a int, b int, amount float64, db *gorm.DB) bool {
	// 开启事务
	db.Begin()

	//查询A中数据，判断转账额度是否满足要求
	var acc1 Account
	data01 := db.Where("id =?", a).Find(&acc1)
	if data01.Error != nil {
		fmt.Println("查询数据异常")
		db.Rollback()
		return false
	}

	if acc1.Balance < amount {
		fmt.Println("查询数据额度太小")
		db.Rollback()
		return false
	}

	//更新a，b用户数据
	aLeft := acc1.Balance - amount
	data02 := db.Model(&Account{}).Where("id = ?", a).Update("balance", aLeft)
	if data02.Error != nil {
		fmt.Println("转出数据异常")
		db.Rollback()
		return false
	}

	acc2 := Account{}
	db.Where("id = ? ", b).Find(&acc2)
	bLeft := acc2.Balance + amount
	data03 := db.Model(&Account{}).Where("id = ?", b).Update("balance", bLeft)
	if data03.Error != nil {
		fmt.Println("转入数据异常")
		db.Rollback()
		return false
	}

	//事务表中增加记录
	cre := Transaction{
		FromAccountId: a,
		ToAccountId:   b,
		Amount:        amount,
	}
	data04 := db.Create(&cre)
	if data04.Error != nil {
		fmt.Println("插入记录数据异常")
		db.Rollback()
		return false
	}

	fmt.Println("执行结束")

	//提交事务
	db.Commit()
	return true

}

// 连接数据库方法
func ConnenDb() (db *gorm.DB, err error) {
	//开启数据库连接
	dsn := "root:asdf@tcp(localhost:3306)/go_study?charset=utf8mb4&parseTime=True&loc=Local"
	db1, err1 := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db1, err1
}

type Account struct {
	ID      int
	Name    string
	Balance float64
}

type Transaction struct {
	ID            int
	FromAccountId int
	ToAccountId   int
	Amount        float64
}
