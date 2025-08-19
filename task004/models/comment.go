package models

type Comment struct {
	Base
	Content string `gorm:"type:text;not null"`
	UserID  uint64 `gorm:"not null"`
	User    User   // 添加用户关联
	PostID  uint   `gorm:"not null"`
	Post    Post   // 添加文章关联
}
