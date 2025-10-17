package weibo

import (
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type WeiboTool struct {
	*utools.BaseTool
}

// 微博热搜 API 响应结构
type WeiboResponse struct {
	Ok   int `json:"ok"`
	Data struct {
		Realtime []WeiboHotItem `json:"realtime"`
	} `json:"data"`
}

type WeiboHotItem struct {
	Word            string      `json:"word"`
	WordScheme      string      `json:"word_scheme"`
	Num             int64       `json:"num"`
	RawHot          int64       `json:"raw_hot"`
	Rank            int         `json:"rank"`
	Note            string      `json:"note"`
	Label           string      `json:"label"`
	Category        string      `json:"category"`
	MidUrl          string      `json:"mid_url"`
	Icon            string      `json:"icon"`
	IconWidth       int         `json:"icon_width"`
	IconHeight      int         `json:"icon_height"`
	Ad              int         `json:"ad"`
	IsNew           int         `json:"is_new"`
	IsFei           int         `json:"is_fei"`
	IsHot           int         `json:"is_hot"`
	IsBoom          int         `json:"is_boom"`
	Flag            int         `json:"flag"`
	FlagDesc        string      `json:"flag_desc"`
	Mid             string      `json:"mid"`
	Onboard_time    int64       `json:"onboard_time"`
	Star_name       interface{} `json:"star_name"`
	Topic_flag      int         `json:"topic_flag"`
	Emoticon        string      `json:"emoticon"`
	IsAbtest        int         `json:"is_abtest"`
	Realpos         int         `json:"realpos"`
}

func NewWeiboTool() *WeiboTool {
	metadata := &utools.ToolMetadata{
		Code:        "weibo_hot",
		Name:        "微博热搜",
		Description: "获取微博实时热搜榜单，包含话题、热度、排名等信息",
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
			"exclude_ads": {
				Type:        "boolean",
				Title:       "排除广告",
				Description: "过滤掉广告热搜",
				Default:     true,
			},
			"min_hot": {
				Type:        "integer",
				Title:       "最小热度值",
				Description: "过滤热度低于此值的热搜，0 表示不过滤",
				Default:     0,
				Minimum:     float64Ptr(0),
			},
			"category_filter": {
				Type:        "string",
				Title:       "分类筛选",
				Description: "只显示特定分类的热搜（如：社会、娱乐、科技），留空则不过滤，多个用逗号分隔",
				Default:     "",
			},
			"exclude_keywords": {
				Type:        "string",
				Title:       "排除关键词",
				Description: "排除包含特定关键词的热搜，多个关键词用逗号分隔",
				Default:     "",
			},
			"only_new": {
				Type:        "boolean",
				Title:       "仅显示新话题",
				Description: "只显示标记为「新」的热搜话题",
				Default:     false,
			},
			"only_hot": {
				Type:        "boolean",
				Title:       "仅显示热门话题",
				Description: "只显示标记为「热」的热搜话题",
				Default:     false,
			},
		},
		Required: []string{},
	}

	return &WeiboTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *WeiboTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	// 1. 解析配置
	maxItems := 10
	if val, ok := config["max_items"].(float64); ok {
		maxItems = int(val)
	} else if val, ok := config["max_items"].(int); ok {
		maxItems = val
	}

	excludeAds := true
	if val, ok := config["exclude_ads"].(bool); ok {
		excludeAds = val
	}

	minHot := 0
	if val, ok := config["min_hot"].(float64); ok {
		minHot = int(val)
	} else if val, ok := config["min_hot"].(int); ok {
		minHot = val
	}

	categoryFilter := ""
	if val, ok := config["category_filter"].(string); ok {
		categoryFilter = strings.TrimSpace(val)
	}

	excludeKeywords := ""
	if val, ok := config["exclude_keywords"].(string); ok {
		excludeKeywords = strings.TrimSpace(val)
	}

	onlyNew := false
	if val, ok := config["only_new"].(bool); ok {
		onlyNew = val
	}

	onlyHot := false
	if val, ok := config["only_hot"].(bool); ok {
		onlyHot = val
	}

	// 2. 调用微博 API
	hotItems, err := t.fetchWeiboHotSearch()
	if err != nil {
		return &utools.ExecutionResult{
			Success: false,
			Message: fmt.Sprintf("获取微博热搜失败: %v", err),
		}, fmt.Errorf("获取微博热搜失败: %v", err)
	}

	// 3. 解析过滤条件
	var categoryList []string
	if categoryFilter != "" {
		for _, cat := range strings.Split(categoryFilter, ",") {
			if trimmed := strings.TrimSpace(cat); trimmed != "" {
				categoryList = append(categoryList, strings.ToLower(trimmed))
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

	// 4. 过滤和处理数据
	var filteredItems []map[string]interface{}
	for _, item := range hotItems {
		// 过滤广告
		if excludeAds && item.Ad == 1 {
			continue
		}

		// 热度过滤
		if minHot > 0 && item.Num < int64(minHot) {
			continue
		}

		// 分类过滤
		if len(categoryList) > 0 {
			matched := false
			itemCategory := strings.ToLower(item.Category)
			for _, cat := range categoryList {
				if strings.Contains(itemCategory, cat) {
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
			searchText := strings.ToLower(item.Word + " " + item.Note)
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

		// 仅显示新话题
		if onlyNew && item.IsNew != 1 {
			continue
		}

		// 仅显示热门话题
		if onlyHot && item.IsHot != 1 {
			continue
		}

		// 构建热搜数据
		itemData := map[string]interface{}{
			"rank":     item.Rank + 1, // API 返回的 rank 从 0 开始
			"word":     item.Word,
			"hot":      item.Num,
			"raw_hot":  item.RawHot,
			"category": item.Category,
			"url":      fmt.Sprintf("https://s.weibo.com/weibo?q=%%23%s%%23", item.Word),
			"note":     item.Note,
		}

		// 添加标签信息
		if item.Label != "" {
			itemData["label"] = item.Label
		}

		// 添加图标描述
		if item.FlagDesc != "" {
			itemData["flag_desc"] = item.FlagDesc
		}

		// 添加明星名字（如果有）
		if item.Star_name != nil {
			if starNames, ok := item.Star_name.([]interface{}); ok && len(starNames) > 0 {
				itemData["star_name"] = starNames[0]
			}
		}

		// 添加上榜时间
		if item.Onboard_time > 0 {
			itemData["onboard_time"] = time.Unix(item.Onboard_time, 0).Format("2006-01-02 15:04:05")
		}

		// 标记类型
		flags := []string{}
		if item.IsNew == 1 {
			flags = append(flags, "新")
		}
		if item.IsHot == 1 {
			flags = append(flags, "热")
		}
		if item.IsBoom == 1 {
			flags = append(flags, "爆")
		}
		if item.IsFei == 1 {
			flags = append(flags, "沸")
		}
		if len(flags) > 0 {
			itemData["flags"] = flags
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
		"source": "Weibo Hot Search API",
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("成功获取 %d 条微博热搜", len(filteredItems)),
		Output: map[string]interface{}{
			"response": response,
			"items":    filteredItems,
			"total":    len(filteredItems),
		},
	}, nil
}

// fetchWeiboHotSearch 调用微博 API 获取热搜
func (t *WeiboTool) fetchWeiboHotSearch() ([]WeiboHotItem, error) {
	apiURL := "https://weibo.com/ajax/side/hotSearch"

	// 创建 HTTP 请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头（重要！）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://weibo.com")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误状态码: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析 JSON
	var weiboResp WeiboResponse
	if err := json.Unmarshal(body, &weiboResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if weiboResp.Ok != 1 {
		return nil, fmt.Errorf("API 返回错误: ok=%d", weiboResp.Ok)
	}

	return weiboResp.Data.Realtime, nil
}

func float64Ptr(v float64) *float64 {
	return &v
}

func init() {
	utools.Register(NewWeiboTool())
}
