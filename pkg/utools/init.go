package utools

import (
	"log"
)

// InitTools 初始化工具系统
func InitTools() {
	log.Println("工具系统初始化中...")

	// 统计已注册工具数量
	count := GetRegistry().Count()
	log.Printf("成功注册 %d 个工具\n", count)

	// 列出所有已注册的工具
	tools := GetRegistry().List()
	for _, tool := range tools {
		log.Printf("  - [%s] %s (v%s) - %s\n", tool.Code, tool.Name, tool.Version, tool.Category)
	}

	log.Println("工具系统初始化完成")
}
