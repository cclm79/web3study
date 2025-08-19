package controllers

import (
	"net/http"
	"strconv"
	"testproject/task004/models"
	"testproject/task004/pkg/mysql"
	"testproject/task004/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	type PostInput struct {
		Title   string `json:"title" binding:"required,min=3,max=255"`
		Content string `json:"content" binding:"required,min=1"`
	}

	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendValidationError(c, utils.GetValidationErrors(err))
		return
	}

	// 创建文章
	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID,
	}

	if err := mysql.MysqlDb().Create(&post).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// 返回创建成功的文章
	utils.SendCreated(c, gin.H{
		"id":        post.ID,
		"title":     post.Title,
		"content":   post.Content,
		"userId":    post.UserID,
		"createdAt": post.CreatedAt,
	})
}

// GetAllPosts 获取所有文章
func GetAllPosts(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	var posts []models.Post

	// 获取文章总数
	if err := mysql.MysqlDb().Model(&models.Post{}).Count(&total).Error; err != nil {
		utils.SendInternalServerError(c)
		return
	}

	// 获取文章列表（包含作者信息）
	if err := mysql.MysqlDb().Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username")
	}).Order("created_at desc").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		utils.SendInternalServerError(c)
		return
	}

	// 简化响应数据
	response := make([]gin.H, len(posts))
	for i, post := range posts {
		response[i] = gin.H{
			"id":        post.ID,
			"title":     post.Title,
			"excerpt":   utils.TruncateString(post.Content, 100),
			"author":    post.User.Username,
			"userId":    post.UserID,
			"createdAt": post.CreatedAt,
			"updatedAt": post.UpdatedAt,
		}
	}

	// 返回分页结果
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"limit": limit,
		"posts": response,
	})
}

// GetPostByID 获取单篇文章
func GetPostByID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	if err := mysql.MysqlDb().Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username")
	}).First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendNotFound(c, "Post")
		} else {
			utils.SendInternalServerError(c)
		}
		return
	}

	// 返回文章详情
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"id":        post.ID,
		"title":     post.Title,
		"content":   post.Content,
		"author":    post.User.Username,
		"userId":    post.UserID,
		"createdAt": post.CreatedAt,
		"updatedAt": post.UpdatedAt,
		"comments":  post.Comments,
	})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)
	post, exists := c.Get("post")
	if !exists {
		utils.SendInternalServerError(c)
		return
	}

	postObj := post.(models.Post)

	// 验证用户是否有权限更新
	if postObj.UserID != userID {
		utils.SendForbidden(c, "You are not the owner of this post")
		return
	}

	type UpdateInput struct {
		Title   string `json:"title" binding:"omitempty,min=3,max=255"`
		Content string `json:"content" binding:"omitempty,min=10"`
	}

	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendValidationError(c, utils.GetValidationErrors(err))
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Content != "" {
		updates["content"] = input.Content
	}

	if len(updates) == 0 {
		utils.SendError(c, http.StatusBadRequest, "No valid fields to update")
		return
	}

	// 更新文章
	if err := mysql.MysqlDb().Model(&postObj).Updates(updates).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to update post")
		return
	}

	// 返回更新后的文章
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"id":        postObj.ID,
		"title":     postObj.Title,
		"content":   postObj.Content,
		"updatedAt": postObj.UpdatedAt,
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)
	post, exists := c.Get("post")
	if !exists {
		utils.SendInternalServerError(c)
		return
	}

	postObj := post.(models.Post)

	// 验证用户是否有权限删除
	if postObj.UserID != userID {
		utils.SendForbidden(c, "You are not the owner of this post")
		return
	}

	// 删除文章
	if err := mysql.MysqlDb().Delete(&postObj).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	// 返回删除成功响应
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"message": "Post deleted successfully",
		"id":      postObj.ID,
	})
}

// PostOwnerCheck 文章所有权检查中间件
func PostOwnerCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint64)
		postID := c.Param("id")

		var post models.Post
		if err := mysql.MysqlDb().First(&post, postID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.SendNotFound(c, "Post")
			} else {
				utils.SendInternalServerError(c)
			}
			c.Abort()
			return
		}
		if post.UserID != userID {
			utils.SendNotFound(c, "You are not the owner of this post")
			return
		}

		// 将文章对象存储到上下文中
		c.Set("post", post)
		c.Next()
	}
}

// 在GetPostByID函数中
func GetPostByID2(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	// 修正预加载方式
	if err := mysql.MysqlDb().Preload("User").Preload("Comments.User").
		First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendNotFound(c, "Post")
		} else {
			utils.SendInternalServerError(c)
		}
		return
	}

	// 返回文章详情
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"id":        post.ID,
		"title":     post.Title,
		"content":   post.Content,
		"author":    post.User.Username,
		"userId":    post.UserID,
		"createdAt": post.CreatedAt,
		"updatedAt": post.UpdatedAt,
		"comments":  post.Comments,
	})
}
