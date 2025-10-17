package baidu

import (
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type BaiduHotTool struct {
	*utools.BaseTool
}

// 百度热搜数据结构
type BaiduHotItem struct {
	Index       int    `json:"index"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	HotScore    string `json:"hotScore"`
	ImgURL      string `json:"imgUrl"`
	PicURL      string `json:"pic_url"`
	ArticleURL  string `json:"articleURL"`
	RawURL      string `json:"rawUrl"`
	MiddleImage string `json:"middleImage"`
}

func NewBaiduHotTool() *BaiduHotTool {
	metadata := &utools.ToolMetadata{
		Code:        "baidu_hot",
		Name:        "百度热搜",
		Description: "获取百度实时热点榜单，包含热搜话题、热度值、链接等信息",
		Category:    "news",
		Version:     "1.0.0",
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "完整响应",
				Children: map[string]utools.OutputFieldDef{
					"total": {
						Type:  "integer",
						Label: "热搜总数",
					},
					"items": {
						Type:  "array",
						Label: "热搜列表",
					},
					"source": {
						Type:  "string",
						Label: "数据源",
					},
				},
			},
			"items": {
				Type:  "array",
				Label: "热搜列表（快捷访问）",
			},
			"total": {
				Type:  "integer",
				Label: "热搜总数（快捷访问）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"max_items": {
				Type:        "integer",
				Title:       "最大条目数",
				Description: "获取热搜的最大数量",
				Default:     10,
				Minimum:     float64Ptr(1),
				Maximum:     float64Ptr(50),
			},
			"min_rank": {
				Type:        "integer",
				Title:       "最小排名",
				Description: "只显示排名在此值以内的热搜（1-50），0 表示不限制",
				Default:     0,
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(50),
			},
			"exclude_keywords": {
				Type:        "string",
				Title:       "排除关键词",
				Description: "排除包含特定关键词的热搜，多个关键词用逗号分隔",
				Default:     "",
			},
		},
		Required: []string{},
	}

	return &BaiduHotTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *BaiduHotTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	// 1. 解析配置
	maxItems := 10
	if val, ok := config["max_items"].(float64); ok {
		maxItems = int(val)
	} else if val, ok := config["max_items"].(int); ok {
		maxItems = val
	}

	minRank := 0
	if val, ok := config["min_rank"].(float64); ok {
		minRank = int(val)
	} else if val, ok := config["min_rank"].(int); ok {
		minRank = val
	}

	excludeKeywords := ""
	if val, ok := config["exclude_keywords"].(string); ok {
		excludeKeywords = strings.TrimSpace(val)
	}

	// 2. 解析排除关键词
	var excludeList []string
	if excludeKeywords != "" {
		for _, kw := range strings.Split(excludeKeywords, ",") {
			if trimmed := strings.TrimSpace(kw); trimmed != "" {
				excludeList = append(excludeList, strings.ToLower(trimmed))
			}
		}
	}

	// 3. 调用百度热搜 API
	hotItems, err := t.fetchBaiduHot()
	if err != nil {
		return &utools.ExecutionResult{
			Success: false,
			Message: fmt.Sprintf("获取百度热搜失败: %v", err),
		}, fmt.Errorf("获取百度热搜失败: %v", err)
	}

	// 4. 过滤和处理数据
	var filteredItems []map[string]interface{}
	for _, item := range hotItems {
		// 排名过滤
		if minRank > 0 && item.Index > minRank {
			continue
		}

		// 关键词排除
		if len(excludeList) > 0 {
			shouldExclude := false
			searchText := strings.ToLower(item.Title)
			for _, kw := range excludeList {
				if strings.Contains(searchText, kw) {
					shouldExclude = true
					break
				}
			}
			if shouldExclude {
				continue
			}
		}

		// 构建热搜数据
		itemData := map[string]interface{}{
			"rank":      item.Index,
			"title":     item.Title,
			"url":       item.URL,
			"hot_score": item.HotScore,
		}

		// 添加图片（如果有）
		if item.ImgURL != "" {
			itemData["img_url"] = item.ImgURL
		} else if item.PicURL != "" {
			itemData["img_url"] = item.PicURL
		}

		filteredItems = append(filteredItems, itemData)

		// 达到最大数量，停止
		if len(filteredItems) >= maxItems {
			break
		}
	}

	// 5. 返回结果
	response := map[string]interface{}{
		"total":  len(filteredItems),
		"items":  filteredItems,
		"source": "Baidu Hot Search",
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("成功获取 %d 条百度热搜", len(filteredItems)),
		Output: map[string]interface{}{
			"response": response,
			"items":    filteredItems,
			"total":    len(filteredItems),
		},
	}, nil
}

// fetchBaiduHot 获取百度热搜数据
func (t *BaiduHotTool) fetchBaiduHot() ([]BaiduHotItem, error) {
	// 百度热搜页面
	apiURL := "https://top.baidu.com/api/board?platform=wise&tab=realtime"

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://top.baidu.com/board?tab=realtime")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析 JSON
	var result struct {
		Data struct {
			Cards []struct {
				Content []BaiduHotItem `json:"content"`
			} `json:"cards"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 提取热搜列表
	var hotItems []BaiduHotItem
	if len(result.Data.Cards) > 0 && len(result.Data.Cards[0].Content) > 0 {
		hotItems = result.Data.Cards[0].Content
	}

	// 处理URL（移除转义字符）
	for i := range hotItems {
		if hotItems[i].URL != "" {
			hotItems[i].URL = unescapeURL(hotItems[i].URL)
		}
	}

	return hotItems, nil
}

// unescapeURL 解码URL中的转义字符
func unescapeURL(s string) string {
	s = strings.ReplaceAll(s, "\\u002F", "/")
	s = strings.ReplaceAll(s, "\\", "")
	return s
}

func float64Ptr(v float64) *float64 {
	return &v
}

func init() {
	utools.Register(NewBaiduHotTool())
}
