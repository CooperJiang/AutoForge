package utools

// ToolCategory 工具分类定义
type ToolCategory struct {
	Code        string `json:"code"`        // 分类代码
	Name        string `json:"name"`        // 分类名称
	Description string `json:"description"` // 分类描述
	Icon        string `json:"icon"`        // 分类图标（可选）
}

// 预定义的工具分类常量
const (
	CategoryAI           = "ai"           // AI 工具
	CategoryData         = "data"         // 数据处理
	CategoryNotification = "notification" // 通知
	CategoryNetwork      = "network"      // 网络
	CategoryStorage      = "storage"      // 存储
	CategoryUtility      = "utility"      // 工具
	CategoryMonitoring   = "monitoring"   // 监控
)

// toolCategories 工具分类列表
var toolCategories = []ToolCategory{
	{
		Code:        CategoryAI,
		Name:        "AI",
		Description: "人工智能相关工具",
		Icon:        "🤖",
	},
	{
		Code:        CategoryData,
		Name:        "数据处理",
		Description: "数据转换、处理工具",
		Icon:        "📊",
	},
	{
		Code:        CategoryNotification,
		Name:        "通知",
		Description: "消息通知工具",
		Icon:        "📧",
	},
	{
		Code:        CategoryNetwork,
		Name:        "网络",
		Description: "网络请求工具",
		Icon:        "🌐",
	},
	{
		Code:        CategoryStorage,
		Name:        "存储",
		Description: "文件存储工具",
		Icon:        "💾",
	},
	{
		Code:        CategoryUtility,
		Name:        "工具",
		Description: "实用工具",
		Icon:        "🔧",
	},
	{
		Code:        CategoryMonitoring,
		Name:        "监控",
		Description: "健康检查、监控工具",
		Icon:        "📡",
	},
}

// GetToolCategories 获取所有工具分类
func GetToolCategories() []ToolCategory {
	return toolCategories
}

// GetCategoryName 根据代码获取分类名称
func GetCategoryName(code string) string {
	for _, category := range toolCategories {
		if category.Code == code {
			return category.Name
		}
	}
	return code
}

// IsCategoryValid 验证分类是否有效
func IsCategoryValid(code string) bool {
	for _, category := range toolCategories {
		if category.Code == code {
			return true
		}
	}
	return false
}
