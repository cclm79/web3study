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

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	type CommentInput struct {
		Content string `json:"content" binding:"required,min=1,max=500"`
	}

	var input CommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendValidationError(c, utils.GetValidationErrors(err))
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := mysql.MysqlDb().First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendNotFound(c, "Post")
		} else {
			utils.SendInternalServerError(c)
		}
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: input.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}

	if err := mysql.MysqlDb().Create(&comment).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// 加载用户信息 - 修正后的查询
	var commentWithUser models.Comment
	if err := mysql.MysqlDb().Preload("User").First(&commentWithUser, comment.ID).Error; err != nil {
		utils.SendInternalServerError(c)
		return
	}

	// 返回创建成功的评论
	utils.SendCreated(c, gin.H{
		"id":        commentWithUser.ID,
		"content":   commentWithUser.Content,
		"userId":    commentWithUser.UserID,
		"username":  commentWithUser.User.Username, // 现在可以安全访问
		"postId":    commentWithUser.PostID,
		"createdAt": commentWithUser.CreatedAt,
	})
}

// GetCommentsByPost 获取文章的所有评论
func GetCommentsByPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	var comments []models.Comment

	// 获取评论总数
	if err := mysql.MysqlDb().Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		utils.SendInternalServerError(c)
		return
	}

	// 获取评论列表（包含用户信息）- 修正后的查询
	if err := mysql.MysqlDb().Preload("User").
		Where("post_id = ?", postID).
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&comments).Error; err != nil {
		utils.SendInternalServerError(c)
		return
	}

	// 简化响应数据
	response := make([]gin.H, len(comments))
	for i, comment := range comments {
		response[i] = gin.H{
			"id":        comment.ID,
			"content":   comment.Content,
			"userId":    comment.UserID,
			"username":  comment.User.Username, // 现在可以安全访问
			"createdAt": comment.CreatedAt,
		}
	}

	// 返回分页结果
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"total":    total,
		"page":     page,
		"limit":    limit,
		"comments": response,
	})
}
