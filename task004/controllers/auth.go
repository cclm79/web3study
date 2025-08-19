package controllers

import (
	"net/http"
	"testproject/task004/models"
	"testproject/task004/pkg/mysql"
	"testproject/task004/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register 用户注册
func Register(c *gin.Context) {
	type RegisterInput struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendValidationError(c, utils.GetValidationErrors(err))
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := mysql.MysqlDb().Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		utils.SendError(c, http.StatusConflict, "Username already exists")
		return
	}

	// 检查邮箱是否已存在
	if err := mysql.MysqlDb().Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		utils.SendError(c, http.StatusConflict, "Email already exists")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendInternalServerError(c)
		return
	}

	// 创建用户
	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
	}

	if err := mysql.MysqlDb().Create(&user).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// 返回创建成功的响应（不包含密码）
	utils.SendCreated(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	type LoginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendValidationError(c, utils.GetValidationErrors(err))
		return
	}

	// 查找用户
	var user models.User
	if err := mysql.MysqlDb().Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendUnauthorized(c, "Invalid username or password")
		} else {
			utils.SendInternalServerError(c)
		}
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.SendUnauthorized(c, "Invalid username or password")
		return
	}

	// 生成JWT
	token, err := utils.GenerateToken(user)
	if err != nil {
		utils.SendInternalServerError(c)
		return
	}

	// 返回登录成功响应
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"token": token,
		// "expires": expiresAt,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetUserProfile 获取当前用户信息
func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "User not authenticated")
		return
	}

	var user models.User
	if err := mysql.MysqlDb().First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendNotFound(c, "User")
		} else {
			utils.SendInternalServerError(c)
		}
		return
	}

	// 返回用户信息（不包含敏感信息）
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}
