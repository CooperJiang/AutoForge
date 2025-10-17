package utools

import (
	"fmt"
	"sync"
)

type Registry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

var (
	globalRegistry *Registry
	once           sync.Once
)

func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			tools: make(map[string]Tool),
		}
	})
	return globalRegistry
}

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

	if _, exists := r.tools[metadata.Code]; exists {
		return fmt.Errorf("tool '%s' already registered", metadata.Code)
	}

	r.tools[metadata.Code] = tool
	return nil
}

func (r *Registry) Get(code string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[code]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not found", code)
	}

	return tool, nil
}

func (r *Registry) List() []*ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*ToolMetadata, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool.GetMetadata())
	}

	return result
}

// GetAllTools 获取所有工具的map（用于同步）
func (r *Registry) GetAllTools() map[string]Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 返回副本，避免外部修改
	toolsCopy := make(map[string]Tool, len(r.tools))
	for k, v := range r.tools {
		toolsCopy[k] = v
	}

	return toolsCopy
}

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

func (r *Registry) Unregister(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[code]; !exists {
		return fmt.Errorf("tool '%s' not registered", code)
	}

	delete(r.tools, code)
	return nil
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tools)
}

func Register(tool Tool) error {
	return GetRegistry().Register(tool)
}

func Get(code string) (Tool, error) {
	return GetRegistry().Get(code)
}

func List() []*ToolMetadata {
	return GetRegistry().List()
}
