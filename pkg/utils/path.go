package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// SetValueByPath 根据路径设置值
// 路径格式: "nodes.0.config.prompt" 表示 nodes[0].config["prompt"]
func SetValueByPath(data interface{}, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return fmt.Errorf("路径不能为空")
	}

	return setValueRecursive(data, parts, value)
}

func setValueRecursive(data interface{}, parts []string, value interface{}) error {
	if len(parts) == 0 {
		return nil
	}

	currentPart := parts[0]
	remainingParts := parts[1:]

	switch v := data.(type) {
	case map[string]interface{}:
		if len(remainingParts) == 0 {
			// 到达最后一层，设置值
			v[currentPart] = value
			return nil
		}

		// 继续递归
		nextData, exists := v[currentPart]
		if !exists {
			// 如果不存在，创建新的 map
			v[currentPart] = make(map[string]interface{})
			nextData = v[currentPart]
		}

		return setValueRecursive(nextData, remainingParts, value)

	case []interface{}:
		// 数组索引
		index, err := strconv.Atoi(currentPart)
		if err != nil {
			return fmt.Errorf("无效的数组索引: %s", currentPart)
		}

		if index < 0 || index >= len(v) {
			return fmt.Errorf("数组索引超出范围: %d", index)
		}

		if len(remainingParts) == 0 {
			v[index] = value
			return nil
		}

		return setValueRecursive(v[index], remainingParts, value)

	default:
		return fmt.Errorf("无法在类型 %T 上设置路径: %s", data, currentPart)
	}
}

// GetValueByPath 根据路径获取值
func GetValueByPath(data interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("路径不能为空")
	}

	return getValueRecursive(data, parts)
}

func getValueRecursive(data interface{}, parts []string) (interface{}, error) {
	if len(parts) == 0 {
		return data, nil
	}

	currentPart := parts[0]
	remainingParts := parts[1:]

	switch v := data.(type) {
	case map[string]interface{}:
		nextData, exists := v[currentPart]
		if !exists {
			return nil, fmt.Errorf("键不存在: %s", currentPart)
		}

		if len(remainingParts) == 0 {
			return nextData, nil
		}

		return getValueRecursive(nextData, remainingParts)

	case []interface{}:
		index, err := strconv.Atoi(currentPart)
		if err != nil {
			return nil, fmt.Errorf("无效的数组索引: %s", currentPart)
		}

		if index < 0 || index >= len(v) {
			return nil, fmt.Errorf("数组索引超出范围: %d", index)
		}

		if len(remainingParts) == 0 {
			return v[index], nil
		}

		return getValueRecursive(v[index], remainingParts)

	default:
		return nil, fmt.Errorf("无法在类型 %T 上获取路径: %s", data, currentPart)
	}
}
