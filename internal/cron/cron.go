package cron

import (
	"auto-forge/pkg/logger"

	"github.com/robfig/cron/v3"
)

var cronManager *cron.Cron

// InitCronManager 初始化定时任务管理器
func InitCronManager() {
	// 创建一个支持秒级别的cron管理器
	cronManager = cron.New(cron.WithSeconds())

	// 注册所有定时任务
	registerTasks()

	// 启动定时任务
	cronManager.Start()
	logger.Info("Cron manager started")

	// 初始化任务调度器
	InitTaskScheduler()

	// 初始化工作流调度器
	InitWorkflowScheduler()
}

// registerTasks 注册所有定时任务
func registerTasks() {
	// 在这里注册其他定时任务
	// registerOtherTask()
}

// Stop 停止所有定时任务
func Stop() {
	if cronManager != nil {
		cronManager.Stop()
		logger.Info("Cron manager stopped")
	}

	// 停止任务调度器
	StopTaskScheduler()

	// 停止工作流调度器
	StopWorkflowScheduler()
} 