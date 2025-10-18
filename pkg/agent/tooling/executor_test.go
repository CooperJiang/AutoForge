package tooling

import (
	"auto-forge/pkg/utools"
	"context"
	"errors"
	"testing"
	"time"
)

// MockTool 模拟工具
type MockTool struct {
	shouldFail   bool
	failCount    int
	currentFails int
	delay        time.Duration
}

func (t *MockTool) GetMetadata() *utools.ToolMetadata {
	return &utools.ToolMetadata{
		Code:        "mock_tool",
		Name:        "Mock Tool",
		Description: "A mock tool for testing",
	}
}

func (t *MockTool) GetSchema() *utools.ConfigSchema {
	return &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"input": {
				Type:        "string",
				Title:       "Input",
				Description: "Test input",
			},
		},
	}
}

func (t *MockTool) Validate(config map[string]interface{}) error {
	return nil
}

func (t *MockTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	// 模拟延迟
	if t.delay > 0 {
		select {
		case <-time.After(t.delay):
		case <-ctx.Context.Done():
			return nil, ctx.Context.Err()
		}
	}

	// 模拟失败
	if t.shouldFail {
		if t.currentFails < t.failCount {
			t.currentFails++
			return nil, errors.New("connection timeout")
		}
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: "Mock execution successful",
		Output: map[string]interface{}{
			"result": "success",
		},
	}, nil
}

// TestToolExecutor_Success 测试成功执行
func TestToolExecutor_Success(t *testing.T) {
	config := DefaultExecutionConfig()
	executor := NewToolExecutor(config)

	mockTool := &MockTool{}
	ctx := context.Background()
	args := map[string]interface{}{"input": "test"}

	result := executor.Execute(ctx, mockTool, args)

	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}

	if result.Attempts != 1 {
		t.Errorf("Expected 1 attempt, got: %d", result.Attempts)
	}
}

// TestToolExecutor_Timeout 测试超时
func TestToolExecutor_Timeout(t *testing.T) {
	config := &ExecutionConfig{
		TimeoutSeconds: 1, // 1 秒超时
	}
	executor := NewToolExecutor(config)

	mockTool := &MockTool{
		delay: 2 * time.Second, // 延迟 2 秒
	}
	ctx := context.Background()
	args := map[string]interface{}{"input": "test"}

	result := executor.Execute(ctx, mockTool, args)

	if result.Error == nil {
		t.Error("Expected timeout error, got nil")
	}

	if !errors.Is(result.Error, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded error, got: %v", result.Error)
	}
}

// TestToolExecutor_Retry 测试重试
func TestToolExecutor_Retry(t *testing.T) {
	config := &ExecutionConfig{
		TimeoutSeconds: 10,
		Retry: &RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    100, // 100ms for faster test
			MaxBackoff:        500,
			BackoffMultiplier: 2.0,
			RetryableErrors: []string{
				"timeout",
				"connection",
			},
		},
	}
	executor := NewToolExecutor(config)

	mockTool := &MockTool{
		shouldFail: true,
		failCount:  2, // 前 2 次失败，第 3 次成功
	}
	ctx := context.Background()
	args := map[string]interface{}{"input": "test"}

	result := executor.Execute(ctx, mockTool, args)

	if result.Error != nil {
		t.Errorf("Expected success after retries, got error: %v", result.Error)
	}

	if result.Attempts != 3 {
		t.Errorf("Expected 3 attempts (1 initial + 2 retries), got: %d", result.Attempts)
	}
}

// MockNonRetryableTool 返回不可重试错误的工具
type MockNonRetryableTool struct {
	MockTool
}

func (t *MockNonRetryableTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	return nil, errors.New("invalid parameter") // 参数错误，不应重试
}

// TestToolExecutor_NonRetryableError 测试不可重试的错误
func TestToolExecutor_NonRetryableError(t *testing.T) {
	config := &ExecutionConfig{
		TimeoutSeconds: 10,
		Retry: &RetryConfig{
			MaxRetries:      2,
			InitialBackoff:  100,
			RetryableErrors: []string{"timeout"},
		},
	}
	executor := NewToolExecutor(config)

	// 创建一个返回不可重试错误的工具
	mockTool := &MockNonRetryableTool{}

	ctx := context.Background()
	args := map[string]interface{}{"input": "test"}

	result := executor.Execute(ctx, mockTool, args)

	if result.Error == nil {
		t.Error("Expected error, got nil")
	}

	if result.Attempts != 1 {
		t.Errorf("Expected 1 attempt (no retry for non-retryable error), got: %d", result.Attempts)
	}
}

// TestRetryConfig_IsRetryable 测试错误可重试判断
func TestRetryConfig_IsRetryable(t *testing.T) {
	config := &RetryConfig{
		RetryableErrors: []string{
			"timeout",
			"connection",
			"503",
		},
	}

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Timeout error",
			err:      errors.New("connection timeout"),
			expected: true,
		},
		{
			name:     "503 error",
			err:      errors.New("HTTP 503 Service Unavailable"),
			expected: true,
		},
		{
			name:     "Connection error",
			err:      errors.New("connection refused"),
			expected: true,
		},
		{
			name:     "Invalid parameter",
			err:      errors.New("invalid parameter"),
			expected: false,
		},
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.IsRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

// TestRetryConfig_GetBackoff 测试退避时间计算
func TestRetryConfig_GetBackoff(t *testing.T) {
	config := &RetryConfig{
		InitialBackoff:    1000,  // 1 秒
		MaxBackoff:        10000, // 10 秒
		BackoffMultiplier: 2.0,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},   // 第 1 次重试
		{1, 2 * time.Second},   // 第 2 次重试
		{2, 4 * time.Second},   // 第 3 次重试
		{3, 8 * time.Second},   // 第 4 次重试
		{4, 10 * time.Second},  // 第 5 次重试（达到最大值）
		{10, 10 * time.Second}, // 超过最大值
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := config.GetBackoff(tt.attempt)
			if result != tt.expected {
				t.Errorf("GetBackoff(%d) = %v, want %v", tt.attempt, result, tt.expected)
			}
		})
	}
}
