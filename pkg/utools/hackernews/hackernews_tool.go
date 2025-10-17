package hackernews

import (
	"auto-forge/pkg/utools"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HackerNewsTool struct {
	*utools.BaseTool
}

// Hacker News Story 结构
type HNStory struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Descendants int    `json:"descendants"` // 评论数
	Type        string `json:"type"`
	Text        string `json:"text"`
}

func NewHackerNewsTool() *HackerNewsTool {
	metadata := &utools.ToolMetadata{
		Code:        "hackernews",
		Name:        "Hacker News",
		Description: "获取 Hacker News 热门技术新闻和讨论，支持多种排序方式",
		Category:    "news",
		Version:     "1.0.0",
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "完整响应",
				Children: map[string]utools.OutputFieldDef{
					"total": {
						Type:  "integer",
						Label: "文章总数",
					},
					"items": {
						Type:  "array",
						Label: "文章列表",
					},
					"source": {
						Type:  "string",
						Label: "数据源",
					},
					"sort_by": {
						Type:  "string",
						Label: "排序方式",
					},
				},
			},
			"items": {
				Type:  "array",
				Label: "文章列表（快捷访问）",
			},
			"total": {
				Type:  "integer",
				Label: "文章总数（快捷访问）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"sort_by": {
				Type:        "string",
				Title:       "排序方式",
				Description: "选择文章排序方式",
				Default:     "top",
				Enum:        []interface{}{"top", "new", "best"},
			},
			"max_items": {
				Type:        "integer",
				Title:       "最大条目数",
				Description: "获取文章的最大数量",
				Default:     10,
				Minimum:     float64Ptr(1),
				Maximum:     float64Ptr(30),
			},
			"min_score": {
				Type:        "integer",
				Title:       "最小评分",
				Description: "过滤评分低于此值的文章，0 表示不过滤",
				Default:     0,
				Minimum:     float64Ptr(0),
			},
			"min_comments": {
				Type:        "integer",
				Title:       "最小评论数",
				Description: "过滤评论数少于此值的文章，0 表示不过滤",
				Default:     0,
				Minimum:     float64Ptr(0),
			},
			"hours_ago": {
				Type:        "integer",
				Title:       "时间范围（小时）",
				Description: "只获取 N 小时内的文章，0 表示不限制",
				Default:     0,
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(72),
			},
			"exclude_keywords": {
				Type:        "string",
				Title:       "排除关键词",
				Description: "排除标题包含特定关键词的文章，多个用逗号分隔",
				Default:     "",
			},
		},
		Required: []string{},
	}

	return &HackerNewsTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *HackerNewsTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	// 1. 解析配置
	sortBy := "top"
	if val, ok := config["sort_by"].(string); ok && val != "" {
		sortBy = val
	}

	maxItems := 10
	if val, ok := config["max_items"].(float64); ok {
		maxItems = int(val)
	} else if val, ok := config["max_items"].(int); ok {
		maxItems = val
	}

	minScore := 0
	if val, ok := config["min_score"].(float64); ok {
		minScore = int(val)
	} else if val, ok := config["min_score"].(int); ok {
		minScore = val
	}

	minComments := 0
	if val, ok := config["min_comments"].(float64); ok {
		minComments = int(val)
	} else if val, ok := config["min_comments"].(int); ok {
		minComments = val
	}

	hoursAgo := 0
	if val, ok := config["hours_ago"].(float64); ok {
		hoursAgo = int(val)
	} else if val, ok := config["hours_ago"].(int); ok {
		hoursAgo = val
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

	// 3. 获取Story IDs
	storyIDs, err := t.fetchStoryIDs(sortBy)
	if err != nil {
		return &utools.ExecutionResult{
			Success: false,
			Message: fmt.Sprintf("获取 Hacker News 失败: %v", err),
		}, fmt.Errorf("获取 Hacker News 失败: %v", err)
	}

	// 4. 计算时间截止点
	var cutoffTime time.Time
	if hoursAgo > 0 {
		cutoffTime = time.Now().Add(-time.Duration(hoursAgo) * time.Hour)
	}

	// 5. 获取并过滤文章详情
	var filteredItems []map[string]interface{}
	for _, storyID := range storyIDs {
		if len(filteredItems) >= maxItems {
			break
		}

		story, err := t.fetchStory(storyID)
		if err != nil {
			continue // 跳过获取失败的
		}

		// 评分过滤
		if minScore > 0 && story.Score < minScore {
			continue
		}

		// 评论数过滤
		if minComments > 0 && story.Descendants < minComments {
			continue
		}

		// 时间过滤
		if hoursAgo > 0 {
			storyTime := time.Unix(story.Time, 0)
			if storyTime.Before(cutoffTime) {
				continue
			}
		}

		// 关键词排除
		if len(excludeList) > 0 {
			shouldExclude := false
			searchText := strings.ToLower(story.Title)
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

		// 构建文章数据
		itemData := map[string]interface{}{
			"id":       story.ID,
			"title":    story.Title,
			"url":      story.URL,
			"score":    story.Score,
			"author":   story.By,
			"comments": story.Descendants,
			"time":     time.Unix(story.Time, 0).Format("2006-01-02 15:04:05"),
			"hn_url":   fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ID),
		}

		// 如果没有外链，可能是 Ask HN 或 Show HN
		if story.URL == "" {
			itemData["url"] = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ID)
			if story.Text != "" {
				itemData["text"] = story.Text
			}
		}

		filteredItems = append(filteredItems, itemData)
	}

	// 6. 返回结果
	response := map[string]interface{}{
		"total":   len(filteredItems),
		"items":   filteredItems,
		"source":  "Hacker News Official API",
		"sort_by": sortBy,
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("成功获取 %d 篇 Hacker News 文章", len(filteredItems)),
		Output: map[string]interface{}{
			"response": response,
			"items":    filteredItems,
			"total":    len(filteredItems),
		},
	}, nil
}

// fetchStoryIDs 获取Story ID列表
func (t *HackerNewsTool) fetchStoryIDs(sortBy string) ([]int, error) {
	var apiURL string
	switch sortBy {
	case "new":
		apiURL = "https://hacker-news.firebaseio.com/v0/newstories.json"
	case "best":
		apiURL = "https://hacker-news.firebaseio.com/v0/beststories.json"
	default: // "top"
		apiURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
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

	var storyIDs []int
	if err := json.Unmarshal(body, &storyIDs); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return storyIDs, nil
}

// fetchStory 获取单个Story详情
func (t *HackerNewsTool) fetchStory(storyID int) (*HNStory, error) {
	apiURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyID)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var story HNStory
	if err := json.Unmarshal(body, &story); err != nil {
		return nil, err
	}

	return &story, nil
}

func float64Ptr(v float64) *float64 {
	return &v
}

func init() {
	utools.Register(NewHackerNewsTool())
}
