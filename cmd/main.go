package main

import (
	"auto-forge/internal/cron"
	"auto-forge/internal/middleware"
	"auto-forge/internal/routes"
	taskService "auto-forge/internal/services/task"
	toolService "auto-forge/internal/services/tool"
	uploadService "auto-forge/internal/services/upload"
	"auto-forge/internal/services/user"
	"auto-forge/pkg/cache"
	"auto-forge/pkg/config"
	"auto-forge/pkg/database"
	"auto-forge/pkg/email"
	"auto-forge/pkg/errors"
	"auto-forge/pkg/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	// 导入工具包以触发工具注册
	_ "auto-forge/pkg/utools/email"
	_ "auto-forge/pkg/utools/feishu"
	_ "auto-forge/pkg/utools/formatter"
	_ "auto-forge/pkg/utools/context"
	_ "auto-forge/pkg/utools/health"
	_ "auto-forge/pkg/utools/http"
	_ "auto-forge/pkg/utools/jsontransform"
	_ "auto-forge/pkg/utools/openai"
	_ "auto-forge/pkg/utools/web"
)

// 应用版本号
const appVersion = "1.0.0"

func main() {
	// 设置时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	// 初始化各种服务
	logger.Init()
	config.InitConfig()
	database.InitDB()
	defer database.Close()
	cache.InitCache()
	email.Init()

	// 初始化基础服务
	user.InitUserService()
	uploadService.InitUploadService()

	// 初始化工具服务
	toolService.InitToolService()

	// 初始化任务服务（必须在 cron 之前）
	taskService.InitTaskService()

	// 初始化定时任务
	cron.InitCronManager()
	defer cron.Stop()

	// 设置 gin 模式
	gin.SetMode(config.GetConfig().App.Mode)

	// 使用 gin.New() 替代 gin.Default() 以便完全控制中间件
	r := gin.New()

	// 添加 recovery 和 logger 中间件
	r.Use(gin.Recovery())

	// 添加错误处理和请求ID中间件
	r.Use(errors.ErrorHandler())

	// 添加 CORS 中间件
	r.Use(middleware.CORSMiddleware())

	r.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务器
	logger.Info("服务启动成功，监听端口: %d，版本: %s", config.GetConfig().App.Port, appVersion)

	if err := r.Run(fmt.Sprintf(":%d", config.GetConfig().App.Port)); err != nil {
		panic(err)
	}
}
