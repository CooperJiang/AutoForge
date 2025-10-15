package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// GenerateWorkflowAPIKey 生成工作流 API Key
func GenerateWorkflowAPIKey() (string, error) {
	randomStr, err := GenerateRandomString(48)
	if err != nil {
		return "", err
	}
	return "wf_" + randomStr, nil
}
