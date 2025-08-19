package models

type Post struct {
	Base
	Title    string    `gorm:"not null"`
	Content  string    `gorm:"type:text;not null"`
	UserID   uint64    `gorm:"not null"`
	User     User      // 添加用户关联
	Comments []Comment // 添加评论关联
}
