package main

import (
	"log"

	"github.com/ResertCursorService/internal/api/handlers"
	"github.com/ResertCursorService/internal/api/router"
	"github.com/ResertCursorService/internal/config"
	repo "github.com/ResertCursorService/internal/repository/postgres"
	"github.com/ResertCursorService/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.Init()

	// 初始化数据库连接
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 初始化仓储层
	adminRepo := repo.NewAdminRepository(db)
	activationRepo := repo.NewActivationCodeRepository(db)

	// 初始化服务层
	adminService := service.NewAdminService(adminRepo)
	activationService := service.NewActivationService(activationRepo)

	// 创建默认管理员账户
	if err := adminService.CreateDefaultAdmin(); err != nil {
		log.Fatal("Failed to create default admin:", err)
	}

	// 初始化处理器
	adminHandler := handlers.NewAdminHandler(adminService)
	activationHandler := handlers.NewActivationHandler(activationService)
	appHandler := handlers.NewAppHandler(activationService)

	// 设置路由
	engine := gin.Default()
	r := router.NewRouter(adminHandler, activationHandler, appHandler)
	r.Setup(engine)

	// 启动服务器
	log.Printf("Server is running on port %s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
