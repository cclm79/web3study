package models

type User struct {
	Base
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Email    string    `gorm:"unique;not null"`
	Posts    []Post    // 添加文章关联
	Comments []Comment // 添加评论关联
}
