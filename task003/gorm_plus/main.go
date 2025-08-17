package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
1. 题目1：模型定义
  - 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
    - 要求 ：
      - 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
      - 编写Go代码，使用Gorm创建这些模型对应的数据库表。
2. 题目2：关联查询
  - 基于上述博客系统的模型定义。
    - 要求 ：
      - 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
      - 编写Go代码，使用Gorm查询评论数量最多的文章信息。
3. 题目3：钩子函数
  - 继续使用博客系统的模型。
    - 要求 ：
      - 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
      - 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

*/

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string
	PostCount int    // 用户文章数量统计字段
	Posts     []Post // 一对多关系: 用户->文章
}

type Post struct {
	ID            uint `gorm:"primaryKey"`
	Title         string
	UserID        uint // 外键
	Content       string
	Comments      []Comment `gorm:"foreignKey:PostID"` // 一对多关系: 文章->评论
	CommentCount  int       // 文章评论数量统计字段
	CommentStatus string    // 评论状态
}

type Comment struct {
	ID      uint `gorm:"primaryKey"`
	Content string
	UserID  uint
	PostID  uint // 外键
}

func main() {

	//连接数据库
	db, err := ConnenDb()
	if err != nil {
		fmt.Println("连接数据库异常")
		return
	}

	//创建表
	// db.AutoMigrate(&User{}, &Post{}, &Comment{})
	// 插入数据
	// initData(db)

	// 查询文章相关信息
	var userId uint = 2
	postList, _ := getUserPostsWithComments(db, userId)
	fmt.Printf("用户id为%d的文章及评论\n", userId)
	for _, item := range postList {
		fmt.Println(item)
	}
	//评论最多
	mostComment, _ := getMostCommentedPost(db)
	fmt.Print("评论最多的文章：")
	fmt.Println(mostComment)

	fmt.Println("=================================")
}

// 连接数据库方法
func ConnenDb() (db *gorm.DB, err error) {
	//开启数据库连接
	dsn := "root:asdf@tcp(localhost:3306)/go_study?charset=utf8mb4&parseTime=True&loc=Local"
	db1, err1 := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db1, err1
}

// 文章查询函数
func getUserPostsWithComments(db *gorm.DB, userID uint) ([]Post, error) {
	var posts []Post
	err := db.
		Preload("Comments"). // 预加载评论
		Where("user_id = ?", userID).
		Find(&posts).Error
	return posts, err
}

// 评论最多
func getMostCommentedPost(db *gorm.DB) (Post, error) {
	var post Post
	err := db.
		Order("comment_count DESC"). // 按评论数降序
		First(&post).Error           // 取第一条
	return post, err
}

// post 创建钩子
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 更新相应用户的文章计数
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1")).
		Error
}

// 评论删除钩子
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var count int64
	// 检查当前文章的评论数量
	if err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count).Error; err != nil {
		return err
	}

	// 如果评论数为0，更新文章状态
	if count == 0 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Updates(map[string]interface{}{
				"comment_status": "无评论",
				"comment_count":  0,
			}).Error
	}

	// 更新评论计数
	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_count", count).
		Error
}

// 创建表插入数据
func initData(db *gorm.DB) {

	//插入数据
	db.Create(GetSampleUsers())
	db.Create(GetSamplePosts())
	db.Create(GetSampleComments())
}

func GetSampleUsers() []User {
	return []User{
		{
			ID:    1,
			Name:  "tech_guru",
			Email: "guru@tech.com",
		},
		{
			ID:    2,
			Name:  "code_master",
			Email: "master@code.dev",
		},
		{
			ID:    3,
			Name:  "dev_learner",
			Email: "learner@dev.net",
		},
	}
}

func GetSamplePosts() []Post {
	return []Post{
		{
			ID:      1,
			Title:   "深入理解Golang并发模型",
			Content: "Go语言的并发模型是其最强大的特性之一...",
			UserID:  1,
		},
		{
			ID:      2,
			Title:   "GORM高级技巧大全",
			Content: "本文将介绍GORM的各种高级用法和最佳实践...",
			UserID:  1,
		},
		{
			ID:      3,
			Title:   "从零构建RESTful API",
			Content: "使用Go和Gin框架构建高性能API服务...",
			UserID:  2,
		},
		{
			ID:      4,
			Title:   "数据库优化实战",
			Content: "如何优化SQL查询提升应用性能...",
			UserID:  2,
		},
		{
			ID:      5,
			Title:   "微服务架构设计模式",
			Content: "微服务架构的常见模式和反模式...",
			UserID:  3,
		},
	}
}

func GetSampleComments() []Comment {
	return []Comment{
		{
			ID:      1,
			Content: "非常有深度的文章！",
			UserID:  2,
			PostID:  1,
		},
		{
			ID:      2,
			Content: "期待更多关于channel的内容",
			UserID:  3,
			PostID:  1,
		},
		{
			ID:      3,
			Content: "GORM的关联查询确实很方便",
			UserID:  1,
			PostID:  2,
		},
		{
			ID:      4,
			Content: "解决了我的实际问题",
			UserID:  3,
			PostID:  3,
		},
		{
			ID:      5,
			Content: "优化后性能提升明显",
			UserID:  1,
			PostID:  4,
		},
		{
			ID:      6,
			Content: "实例代码能否分享一下？",
			UserID:  3,
			PostID:  4,
		},
		{
			ID:      7,
			Content: "架构设计思路很清晰",
			UserID:  2,
			PostID:  5,
		},
	}
}
