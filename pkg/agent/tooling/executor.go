package tooling

import (
	"auto-forge/pkg/utools"
	"context"
	"fmt"
	"time"
)

// ToolExecutor 工具执行器（支持超时、重试、缓存）
type ToolExecutor struct {
	config       *ExecutionConfig
	cacheManager *CacheManager
}

// NewToolExecutor 创建工具执行器
func NewToolExecutor(config *ExecutionConfig) *ToolExecutor {
	if config == nil {
		config = DefaultExecutionConfig()
	}
	return &ToolExecutor{
		config:       config,
		cacheManager: NewCacheManager(),
	}
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	Output    interface{}
	Error     error
	Attempts  int           // 尝试次数
	Duration  time.Duration // 总耗时
	FromCache bool          // 是否来自缓存
}

// Execute 执行工具（带超时、重试、缓存）
func (e *ToolExecutor) Execute(
	ctx context.Context,
	tool utools.Tool,
	args map[string]interface{},
) *ExecutionResult {
	startTime := time.Now()
	result := &ExecutionResult{}

	// 1. 检查缓存
	if e.config.Cache != nil && e.config.Cache.Enabled {
		cacheKey := e.cacheManager.GenerateCacheKey(tool.GetMetadata().Code, args)
		if cached, found := e.cacheManager.Get(cacheKey); found {
			result.Output = cached
			result.FromCache = true
			result.Attempts = 0
			result.Duration = time.Since(startTime)
			return result
		}
	}

	// 2. 应用超时
	timeout := time.Duration(e.config.TimeoutSeconds) * time.Second
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// 执行（带重试）
	var lastErr error
	maxAttempts := 1
	if e.config.Retry != nil {
		maxAttempts = e.config.Retry.MaxRetries + 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		result.Attempts = attempt + 1

		// 如果不是第一次尝试，先等待退避时间
		if attempt > 0 && e.config.Retry != nil {
			backoff := e.config.Retry.GetBackoff(attempt - 1)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				result.Error = ctx.Err()
				result.Duration = time.Since(startTime)
				return result
			}
		}

		// 执行工具
		execCtx := &utools.ExecutionContext{
			Context: ctx,
		}

		output, err := tool.Execute(execCtx, args)

		// 成功
		if err == nil {
			result.Output = output
			result.Duration = time.Since(startTime)

			// 3. 存入缓存
			if ShouldCache(e.config.Cache, output, nil) {
				cacheKey := e.cacheManager.GenerateCacheKey(tool.GetMetadata().Code, args)
				cacheTTL := GetCacheTTL(e.config.Cache)
				if err := e.cacheManager.Set(cacheKey, output, cacheTTL); err != nil {
					// 缓存失败不影响结果
					fmt.Printf("缓存设置失败: %v\n", err)
				}
			}

			return result
		}

		// 失败
		lastErr = err

		// 检查是否可重试
		if e.config.Retry == nil || !e.config.Retry.IsRetryable(err) {
			break
		}

		// 检查上下文是否已取消
		if ctx.Err() != nil {
			break
		}
	}

	// 所有尝试都失败
	result.Error = fmt.Errorf("工具执行失败（尝试 %d 次）: %w", result.Attempts, lastErr)
	result.Duration = time.Since(startTime)
	return result
}

// ExecuteWithProgress 执行工具并报告进度
func (e *ToolExecutor) ExecuteWithProgress(
	ctx context.Context,
	tool utools.Tool,
	args map[string]interface{},
	progressCallback func(attempt int, message string),
) *ExecutionResult {
	startTime := time.Now()
	result := &ExecutionResult{}

	// 应用超时
	timeout := time.Duration(e.config.TimeoutSeconds) * time.Second
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// 执行（带重试）
	var lastErr error
	maxAttempts := 1
	if e.config.Retry != nil {
		maxAttempts = e.config.Retry.MaxRetries + 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		result.Attempts = attempt + 1

		// 报告进度
		if progressCallback != nil {
			if attempt == 0 {
				progressCallback(attempt+1, "开始执行")
			} else {
				progressCallback(attempt+1, fmt.Sprintf("重试中（第 %d 次）", attempt))
			}
		}

		// 如果不是第一次尝试，先等待退避时间
		if attempt > 0 && e.config.Retry != nil {
			backoff := e.config.Retry.GetBackoff(attempt - 1)

			if progressCallback != nil {
				progressCallback(attempt+1, fmt.Sprintf("等待 %v 后重试", backoff))
			}

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				result.Error = ctx.Err()
				result.Duration = time.Since(startTime)
				return result
			}
		}

		// 执行工具
		execCtx := &utools.ExecutionContext{
			Context: ctx,
		}

		output, err := tool.Execute(execCtx, args)

		// 成功
		if err == nil {
			result.Output = output
			result.Duration = time.Since(startTime)

			if progressCallback != nil {
				progressCallback(attempt+1, "执行成功")
			}

			return result
		}

		// 失败
		lastErr = err

		if progressCallback != nil {
			progressCallback(attempt+1, fmt.Sprintf("执行失败: %v", err))
		}

		// 检查是否可重试
		if e.config.Retry == nil || !e.config.Retry.IsRetryable(err) {
			break
		}

		// 检查上下文是否已取消
		if ctx.Err() != nil {
			break
		}
	}

	// 所有尝试都失败
	result.Error = fmt.Errorf("工具执行失败（尝试 %d 次）: %w", result.Attempts, lastErr)
	result.Duration = time.Since(startTime)
	return result
}
