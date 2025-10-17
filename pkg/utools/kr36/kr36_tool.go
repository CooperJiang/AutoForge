package kr36

import (
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type KR36Tool struct {
	*utools.BaseTool
}

// 36氪快讯数据结构
type KR36Item struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	NewsURL     string `json:"news_url"`
	PublishedAt string `json:"published_at"`
	CoverURL    string `json:"cover_url"`
	ItemType    string `json:"item_type"`
}

type KR36Response struct {
	Code int `json:"code"`
	Data struct {
		Items []struct {
			ItemID      int    `json:"item_id"`
			TemplateType string `json:"template_type"`
			Data        struct {
				ID           int    `json:"id"`
				Title        string `json:"title"`
				Summary      string `json:"summary"`
				NewsURL      string `json:"news_url"`
				PublishedAt  string `json:"published_at"`
				CoverURL     string `json:"cover"`
			} `json:"data"`
		} `json:"items"`
	} `json:"data"`
}

func NewKR36Tool() *KR36Tool {
	metadata := &utools.ToolMetadata{
		Code:        "kr36_news",
		Name:        "36氪快讯",
		Description: "获取36氪最新科技创投快讯，聚焦创业公司和投资动态",
		Category:    "news",
		Version:     "1.0.0",
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "完整响应",
				Children: map[string]utools.OutputFieldDef{
					"total": {
						Type:  "integer",
						Label: "快讯总数",
					},
					"items": {
						Type:  "array",
						Label: "快讯列表",
					},
					"source": {
						Type:  "string",
						Label: "数据源",
					},
				},
			},
			"items": {
				Type:  "array",
				Label: "快讯列表（快捷访问）",
			},
			"total": {
				Type:  "integer",
				Label: "快讯总数（快捷访问）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"max_items": {
				Type:        "integer",
				Title:       "最大条目数",
				Description: "获取快讯的最大数量",
				Default:     10,
				Minimum:     float64Ptr(1),
				Maximum:     float64Ptr(30),
			},
			"hours_ago": {
				Type:        "integer",
				Title:       "时间范围（小时）",
				Description: "只获取 N 小时内的快讯，0 表示不限制",
				Default:     0,
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(72),
			},
			"keywords": {
				Type:        "string",
				Title:       "关键词筛选",
				Description: "只显示包含特定关键词的快讯，多个关键词用逗号分隔（满足任一即可）",
				Default:     "",
			},
			"exclude_keywords": {
				Type:        "string",
				Title:       "排除关键词",
				Description: "排除包含特定关键词的快讯，多个关键词用逗号分隔",
				Default:     "",
			},
		},
		Required: []string{},
	}

	return &KR36Tool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *KR36Tool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	// 1. 解析配置
	maxItems := 10
	if val, ok := config["max_items"].(float64); ok {
		maxItems = int(val)
	} else if val, ok := config["max_items"].(int); ok {
		maxItems = val
	}

	hoursAgo := 0
	if val, ok := config["hours_ago"].(float64); ok {
		hoursAgo = int(val)
	} else if val, ok := config["hours_ago"].(int); ok {
		hoursAgo = val
	}

	keywords := ""
	if val, ok := config["keywords"].(string); ok {
		keywords = strings.TrimSpace(val)
	}

	excludeKeywords := ""
	if val, ok := config["exclude_keywords"].(string); ok {
		excludeKeywords = strings.TrimSpace(val)
	}

	// 2. 解析关键词
	var keywordList []string
	if keywords != "" {
		for _, kw := range strings.Split(keywords, ",") {
			if trimmed := strings.TrimSpace(kw); trimmed != "" {
				keywordList = append(keywordList, strings.ToLower(trimmed))
			}
		}
	}

	var excludeList []string
	if excludeKeywords != "" {
		for _, kw := range strings.Split(excludeKeywords, ",") {
			if trimmed := strings.TrimSpace(kw); trimmed != "" {
				excludeList = append(excludeList, strings.ToLower(trimmed))
			}
		}
	}

	// 3. 获取36氪快讯
	newsItems, err := t.fetchKR36News()
	if err != nil {
		return &utools.ExecutionResult{
			Success: false,
			Message: fmt.Sprintf("获取36氪快讯失败: %v", err),
		}, fmt.Errorf("获取36氪快讯失败: %v", err)
	}

	// 4. 计算时间截止点
	var cutoffTime time.Time
	if hoursAgo > 0 {
		cutoffTime = time.Now().Add(-time.Duration(hoursAgo) * time.Hour)
	}

	// 5. 过滤和处理数据
	var filteredItems []map[string]interface{}
	for _, item := range newsItems {
		// 时间过滤
		if hoursAgo > 0 && item.PublishedAt != "" {
			publishTime, err := time.Parse("2006-01-02 15:04:05", item.PublishedAt)
			if err == nil && publishTime.Before(cutoffTime) {
				continue
			}
		}

		// 关键词筛选
		if len(keywordList) > 0 {
			matched := false
			searchText := strings.ToLower(item.Title + " " + item.Summary)
			for _, kw := range keywordList {
				if strings.Contains(searchText, kw) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// 关键词排除
		if len(excludeList) > 0 {
			shouldExclude := false
			searchText := strings.ToLower(item.Title + " " + item.Summary)
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

		// 构建快讯数据
		itemData := map[string]interface{}{
			"id":           item.ID,
			"title":        item.Title,
			"summary":      item.Summary,
			"url":          item.NewsURL,
			"published_at": item.PublishedAt,
		}

		if item.CoverURL != "" {
			itemData["cover_url"] = item.CoverURL
		}

		filteredItems = append(filteredItems, itemData)

		// 达到最大数量，停止
		if len(filteredItems) >= maxItems {
			break
		}
	}

	// 6. 返回结果
	response := map[string]interface{}{
		"total":  len(filteredItems),
		"items":  filteredItems,
		"source": "36氪快讯",
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("成功获取 %d 条36氪快讯", len(filteredItems)),
		Output: map[string]interface{}{
			"response": response,
			"items":    filteredItems,
			"total":    len(filteredItems),
		},
	}, nil
}

// fetchKR36News 获取36氪快讯数据
func (t *KR36Tool) fetchKR36News() ([]KR36Item, error) {
	// 36氪快讯 API
	apiURL := "https://36kr.com/api/newsflash"

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
	req.Header.Set("Referer", "https://36kr.com/newsflashes")

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
	var result KR36Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 提取快讯列表
	var newsItems []KR36Item
	for _, item := range result.Data.Items {
		if item.Data.ID == 0 {
			continue
		}

		newsItem := KR36Item{
			ID:          item.Data.ID,
			Title:       item.Data.Title,
			Summary:     item.Data.Summary,
			NewsURL:     item.Data.NewsURL,
			PublishedAt: item.Data.PublishedAt,
			CoverURL:    item.Data.CoverURL,
		}

		// 如果没有 NewsURL，使用快讯详情页
		if newsItem.NewsURL == "" {
			newsItem.NewsURL = fmt.Sprintf("https://36kr.com/newsflashes/%d", newsItem.ID)
		}

		newsItems = append(newsItems, newsItem)
	}

	return newsItems, nil
}

func float64Ptr(v float64) *float64 {
	return &v
}

func init() {
	utools.Register(NewKR36Tool())
}
