package main

import (
	"flag"
	"log"
	"strconv"
	"testproject/task004/config"
	"testproject/task004/models"
	logs "testproject/task004/pkg/logger"
	"testproject/task004/pkg/mysql"
	"testproject/task004/routes"

	"github.com/gin-gonic/gin"
)

var configFile = flag.String("f", "./config.yaml", "")

func main() {
	// 设置运行环境
	// setupEnvironment()

	// 初始化配置文件
	config.InitAppConfig(*configFile)

	// 自动迁移数据库模型
	migrateDatabase()

	// 设置路由
	server := setupRouter()

	// 启动服务器
	startServer(server)
}

// migrateDatabase 自动迁移数据库
func migrateDatabase() {
	// 在开发环境自动迁移，生产环境需要谨慎
	if config.Server.RunMode == "debug" {
		log.Println("Starting database migration...")
		mysql.AutoMigrate(
			&models.User{},
			&models.Post{},
			&models.Comment{},
		)
	} else {
		log.Println("Skipping auto migration in production")
	}
}

// setupRouter 设置路由
func setupRouter() *gin.Engine {
	server := gin.Default()

	// 添加中间件
	server.Use(logs.Logger())
	server.Use(gin.Recovery())

	// 设置路由
	routes.SetupRoutes(server)

	return server
}

// startServer 启动服务器
func startServer(router *gin.Engine) {
	port := config.Server.Port
	if port == 0 {
		port = 8080
	}

	log.Printf("Starting server on port %s in %s mode", port, config.Server.RunMode)
	log.Printf("API Documentation: http://localhost:%s/swagger/index.html (if enabled)", port)

	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
