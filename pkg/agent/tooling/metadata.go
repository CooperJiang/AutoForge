package tooling

import "time"

// ExecutionConfig 工具执行配置
type ExecutionConfig struct {
	// 超时配置
	TimeoutSeconds int `json:"timeout_seconds,omitempty"` // 工具级别超时（0 表示使用默认值）

	// 重试配置
	Retry *RetryConfig `json:"retry,omitempty"`

	// 依赖配置
	Dependencies *DependencyConfig `json:"dependencies,omitempty"`

	// 缓存配置
	Cache *CacheConfig `json:"cache,omitempty"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries        int      `json:"max_retries"`        // 最大重试次数
	InitialBackoff    int      `json:"initial_backoff"`    // 初始退避时间（毫秒）
	MaxBackoff        int      `json:"max_backoff"`        // 最大退避时间（毫秒）
	BackoffMultiplier float64  `json:"backoff_multiplier"` // 退避倍数
	RetryableErrors   []string `json:"retryable_errors"`   // 可重试的错误关键词
}

// DependencyConfig 依赖配置
type DependencyConfig struct {
	// 输入类型要求（如 "image_url", "file_object"）
	Requires []string `json:"requires,omitempty"`

	// 输出类型提供（如 "cdn_url", "qrcode_data"）
	Provides []string `json:"provides,omitempty"`

	// 建议的前置工具
	SuggestedPredecessors []string `json:"suggested_predecessors,omitempty"`

	// 互斥工具（不能同时使用）
	ConflictsWith []string `json:"conflicts_with,omitempty"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled bool          `json:"enabled"`       // 是否启用缓存
	TTL     time.Duration `json:"ttl"`           // 缓存过期时间
	Key     string        `json:"key,omitempty"` // 自定义缓存 key 模板
}

// DefaultExecutionConfig 返回默认执行配置
func DefaultExecutionConfig() *ExecutionConfig {
	return &ExecutionConfig{
		TimeoutSeconds: 300, // 默认 5 分钟
		Retry: &RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    1000,  // 1 秒
			MaxBackoff:        10000, // 10 秒
			BackoffMultiplier: 2.0,
			RetryableErrors: []string{
				"timeout",
				"connection",
				"network",
				"rate limit",
				"503",
				"504",
			},
		},
		Cache: &CacheConfig{
			Enabled: false,
			TTL:     5 * time.Minute,
		},
	}
}

// IsRetryable 判断错误是否可重试
func (r *RetryConfig) IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	for _, keyword := range r.RetryableErrors {
		if contains(errMsg, keyword) {
			return true
		}
	}
	return false
}

// GetBackoff 计算退避时间
func (r *RetryConfig) GetBackoff(attempt int) time.Duration {
	backoff := float64(r.InitialBackoff)
	for i := 0; i < attempt; i++ {
		backoff *= r.BackoffMultiplier
		if int(backoff) > r.MaxBackoff {
			backoff = float64(r.MaxBackoff)
			break
		}
	}
	return time.Duration(backoff) * time.Millisecond
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
