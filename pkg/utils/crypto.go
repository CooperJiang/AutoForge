package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"

	"auto-forge/pkg/config"
)

// getEncryptionKey 获取加密密钥
func getEncryptionKey() []byte {
	// 优先从环境变量获取
	key := os.Getenv("TOOL_CONFIG_ENCRYPTION_KEY")
	if key == "" {
		// 如果未设置，从 JWT secret 派生
		key = config.GetConfig().JWT.SecretKey
	}

	// 使用 SHA-256 生成固定长度的密钥（32字节）
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// EncryptToolConfig 加密工具配置
func EncryptToolConfig(configMap map[string]interface{}) (string, error) {
	// 将 map 转换为 JSON
	jsonData, err := json.Marshal(configMap)
	if err != nil {
		return "", err
	}

	// 获取加密密钥
	key := getEncryptionKey()

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)

	// Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptToolConfig 解密工具配置
func DecryptToolConfig(encryptedData string) (map[string]interface{}, error) {
	if encryptedData == "" {
		return make(map[string]interface{}), nil
	}

	// Base64 解码
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// 获取加密密钥
	key := getEncryptionKey()

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 提取 nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	var configMap map[string]interface{}
	if err := json.Unmarshal(plaintext, &configMap); err != nil {
		return nil, err
	}

	return configMap, nil
}
