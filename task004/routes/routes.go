package routes

import (
	"testproject/task004/controllers"
	"testproject/task004/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(router *gin.Engine) *gin.Engine {
	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 认证路由组
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/profile", middlewares.JWTAuth(), controllers.GetUserProfile)
	}

	// API路由组（需要认证）
	api := router.Group("/api")
	api.Use(middlewares.JWTAuth())
	{
		// 文章路由
		setupPostRoutes(api)

		// 评论路由
		setupCommentRoutes(api)
	}

	return router
}

// setupPostRoutes 配置文章相关路由
func setupPostRoutes(api *gin.RouterGroup) {
	// 不需要所有权的操作
	api.GET("/posts", controllers.GetAllPosts)
	api.GET("/posts/:id", controllers.GetPostByID)
	api.POST("/posts", controllers.CreatePost)

	// 需要所有权的操作
	ownerRoutes := api.Group("")
	ownerRoutes.Use(controllers.PostOwnerCheck()) // 检查文章所有权
	{
		ownerRoutes.PUT("/posts/:id", controllers.UpdatePost)
		ownerRoutes.DELETE("/posts/:id", controllers.DeletePost)
	}
}

// setupCommentRoutes 配置评论相关路由
func setupCommentRoutes(api *gin.RouterGroup) {
	api.POST("/posts/:id/comments", controllers.CreateComment)
	api.GET("/posts/:id/comments", controllers.GetCommentsByPost)

	// 可选：评论删除功能
	// api.DELETE("/comments/:commentId", controllers.DeleteComment)
}
