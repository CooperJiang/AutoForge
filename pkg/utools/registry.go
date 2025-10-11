package utools

import (
	"fmt"
	"sync"
)

// Registry 工具注册表 - 全局单例
type Registry struct {
	tools map[string]Tool // code -> tool
	mu    sync.RWMutex
}

var (
	globalRegistry *Registry
	once           sync.Once
)

// GetRegistry 获取全局工具注册表
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			tools: make(map[string]Tool),
		}
	})
	return globalRegistry
}

// Register 注册工具
func (r *Registry) Register(tool Tool) error {
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}

	metadata := tool.GetMetadata()
	if metadata == nil {
		return fmt.Errorf("tool metadata cannot be nil")
	}

	if metadata.Code == "" {
		return fmt.Errorf("tool code cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查是否已注册
	if _, exists := r.tools[metadata.Code]; exists {
		return fmt.Errorf("tool '%s' already registered", metadata.Code)
	}

	r.tools[metadata.Code] = tool
	return nil
}

// Get 根据 code 获取工具
func (r *Registry) Get(code string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[code]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not found", code)
	}

	return tool, nil
}

// List 列出所有已注册的工具
func (r *Registry) List() []*ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*ToolMetadata, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool.GetMetadata())
	}

	return result
}

// ListByCategory 列出指定分类的工具
func (r *Registry) ListByCategory(category string) []*ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*ToolMetadata, 0)
	for _, tool := range r.tools {
		metadata := tool.GetMetadata()
		if metadata.Category == category {
			result = append(result, metadata)
		}
	}

	return result
}

// ListAICallable 列出所有可被 AI 调用的工具
func (r *Registry) ListAICallable() []*ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*ToolMetadata, 0)
	for _, tool := range r.tools {
		metadata := tool.GetMetadata()
		if metadata.AICallable {
			result = append(result, metadata)
		}
	}

	return result
}

// Unregister 取消注册工具
func (r *Registry) Unregister(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[code]; !exists {
		return fmt.Errorf("tool '%s' not registered", code)
	}

	delete(r.tools, code)
	return nil
}

// Count 获取已注册工具数量
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tools)
}

// Register 全局注册函数（简便方法）
func Register(tool Tool) error {
	return GetRegistry().Register(tool)
}

// Get 全局获取函数（简便方法）
func Get(code string) (Tool, error) {
	return GetRegistry().Get(code)
}

// List 全局列表函数（简便方法）
func List() []*ToolMetadata {
	return GetRegistry().List()
}
