package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 基础响应结构
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// SendSuccess 发送成功响应
func SendSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Success:   true,
		Message:   "Operation successful",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// SendCreated 发送资源创建成功响应
func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success:   true,
		Message:   "Resource created successfully",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// SendError 发送错误响应
func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().Unix(),
	})
	c.Abort()
}

// SendValidationError 发送表单验证错误响应
func SendValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, Response{
		Success:   false,
		Error:     "Validation failed",
		Data:      errors,
		Timestamp: time.Now().Unix(),
	})
	c.Abort()
}

// SendUnauthorized 发送未授权错误响应
func SendUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().Unix(),
	})
	c.Abort()
}

// SendForbidden 发送禁止访问错误响应
func SendForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden resource"
	}
	c.JSON(http.StatusForbidden, Response{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().Unix(),
	})
	c.Abort()
}

// SendNotFound 发送资源未找到错误响应
func SendNotFound(c *gin.Context, resource string) {
	message := "Resource not found"
	if resource != "" {
		message = resource + " not found"
	}
	c.JSON(http.StatusNotFound, Response{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().Unix(),
	})
	c.Abort()
}

// SendInternalServerError 发送服务器内部错误响应
func SendInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Success:   false,
		Error:     "Internal server error",
		Timestamp: time.Now().Unix(),
	})
	c.Abort()
}
