package tooling

import (
	"auto-forge/pkg/cache"
	"auto-forge/pkg/logger"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// CacheManager 工具缓存管理器
type CacheManager struct {
	cache cache.Cache
}

// NewCacheManager 创建缓存管理器
func NewCacheManager() *CacheManager {
	return &CacheManager{
		cache: cache.GetCache(),
	}
}

// GenerateCacheKey 生成缓存 key
func (m *CacheManager) GenerateCacheKey(toolName string, args map[string]interface{}) string {
	// 将参数序列化为 JSON
	argsJSON, err := json.Marshal(args)
	if err != nil {
		logger.Warn("序列化参数失败: %v", err)
		return ""
	}

	// 使用 SHA256 生成哈希
	hash := sha256.Sum256(argsJSON)
	hashStr := fmt.Sprintf("%x", hash)

	// 格式: tool:cache:{tool_name}:{hash}
	return fmt.Sprintf("tool:cache:%s:%s", toolName, hashStr[:16])
}

// Get 从缓存获取结果
func (m *CacheManager) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	value, err := m.cache.Get(key)
	if err != nil {
		return nil, false
	}

	// 反序列化
	var result interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		logger.Warn("反序列化缓存失败: %v", err)
		return nil, false
	}

	logger.Info("缓存命中: %s", key)
	return result, true
}

// Set 设置缓存
func (m *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	if key == "" {
		return fmt.Errorf("cache key is empty")
	}

	// 序列化
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存值失败: %w", err)
	}

	// 存储
	if err := m.cache.Set(key, string(valueJSON), ttl); err != nil {
		return fmt.Errorf("存储缓存失败: %w", err)
	}

	logger.Info("缓存已设置: %s (TTL: %v)", key, ttl)
	return nil
}

// Delete 删除缓存
func (m *CacheManager) Delete(key string) error {
	if key == "" {
		return nil
	}

	return m.cache.Del(key)
}

// Clear 清空工具的所有缓存
func (m *CacheManager) Clear(toolName string) error {
	// 注意：这需要 cache 支持模式匹配删除
	// 如果不支持，可以维护一个工具缓存 key 列表
	pattern := fmt.Sprintf("tool:cache:%s:*", toolName)
	logger.Info("清空工具缓存: %s", pattern)

	// 这里简化处理，实际需要根据 cache 实现调整
	return nil
}

// GetStats 获取缓存统计
func (m *CacheManager) GetStats(toolName string) *CacheStats {
	// 这里返回简单统计，实际可以从 metrics 获取
	return &CacheStats{
		ToolName: toolName,
		Hits:     0,
		Misses:   0,
		HitRate:  0,
	}
}

// CacheStats 缓存统计
type CacheStats struct {
	ToolName string  `json:"tool_name"`
	Hits     int64   `json:"hits"`
	Misses   int64   `json:"misses"`
	HitRate  float64 `json:"hit_rate"`
}

// ShouldCache 判断是否应该缓存
func ShouldCache(config *CacheConfig, result interface{}, err error) bool {
	if config == nil || !config.Enabled {
		return false
	}

	// 只缓存成功的结果
	if err != nil {
		return false
	}

	// 结果不能为空
	if result == nil {
		return false
	}

	return true
}

// GetCacheTTL 获取缓存 TTL
func GetCacheTTL(config *CacheConfig) time.Duration {
	if config == nil || config.TTL <= 0 {
		return 5 * time.Minute // 默认 5 分钟
	}
	return config.TTL
}
