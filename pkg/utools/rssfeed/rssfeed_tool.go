package rssfeed

import (
	"auto-forge/pkg/utools"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type RSSFeedTool struct {
	*utools.BaseTool
}

func NewRSSFeedTool() *RSSFeedTool {
	metadata := &utools.ToolMetadata{
		Code:        "rss_feed",
		Name:        "RSS 订阅采集器",
		Description: "解析 RSS/Atom/JSON Feed，获取最新文章和资讯",
		Category:    "data",
		Version:     "1.0.0",
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "完整响应",
				Children: map[string]utools.OutputFieldDef{
					"total_sources": {
						Type:  "integer",
						Label: "订阅源总数",
					},
					"total_items": {
						Type:  "integer",
						Label: "文章总数（去重后）",
					},
					"sources": {
						Type:  "array",
						Label: "订阅源统计信息",
					},
					"sources_with_items": {
						Type:  "array",
						Label: "分组数据（按订阅源分开）",
					},
					"items": {
						Type:  "array",
						Label: "合并数据（所有文章）",
					},
				},
			},
			"items": {
				Type:  "array",
				Label: "合并数组（快捷访问，所有文章）",
			},
			"sources_with_items": {
				Type:  "array",
				Label: "分组数据（快捷访问，按源分开）",
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
			"sources": {
				Type:        "array",
				Title:       "RSS 订阅源列表",
				Description: "可以添加多个订阅源，支持独立配置关键词过滤",
			},
			"max_items": {
				Type:        "integer",
				Title:       "最大条目数（总计）",
				Description: "所有订阅源汇总后的最大文章数量",
				Default:     20,
				Minimum:     float64Ptr(1),
				Maximum:     float64Ptr(200),
			},
			"hours_ago": {
				Type:        "integer",
				Title:       "时间范围（小时）",
				Description: "只获取 N 小时内的文章，0 表示不限制",
				Default:     0,
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(720), // 30天
			},
			"dedup_by": {
				Type:        "string",
				Title:       "去重规则",
				Description: "按链接或标题去重",
				Default:     "link",
				Enum:        []interface{}{"link", "title"},
			},
			"sort_by": {
				Type:        "string",
				Title:       "排序方式",
				Description: "按发布时间或订阅源顺序排序",
				Default:     "time",
				Enum:        []interface{}{"time", "source"},
			},
		},
		Required: []string{"sources"},
	}

	return &RSSFeedTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *RSSFeedTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	// 1. 解析配置
	sourcesRaw, ok := config["sources"].([]interface{})
	if !ok || len(sourcesRaw) == 0 {
		return &utools.ExecutionResult{
			Success: false,
			Message: "至少需要配置一个 RSS 订阅源",
		}, fmt.Errorf("至少需要配置一个 RSS 订阅源")
	}

	// 解析订阅源列表
	type Source struct {
		URL      string
		Keywords string
	}
	var sources []Source
	for _, s := range sourcesRaw {
		if sourceMap, ok := s.(map[string]interface{}); ok {
			url, _ := sourceMap["url"].(string)
			keywords, _ := sourceMap["keywords"].(string)
			if url != "" {
				sources = append(sources, Source{
					URL:      url,
					Keywords: strings.TrimSpace(keywords),
				})
			}
		}
	}

	if len(sources) == 0 {
		return &utools.ExecutionResult{
			Success: false,
			Message: "没有有效的订阅源",
		}, fmt.Errorf("没有有效的订阅源")
	}

	maxItems := 20
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

	dedupBy := "link"
	if val, ok := config["dedup_by"].(string); ok && val != "" {
		dedupBy = val
	}

	sortBy := "time"
	if val, ok := config["sort_by"].(string); ok && val != "" {
		sortBy = val
	}

	// 2. 循环采集所有订阅源
	// 配置 HTTP 客户端（增加超时时间，禁用 IPv6）
	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Timeout: 30 * time.Second, // 30秒超时
		Transport: &http.Transport{
			// 强制使用 IPv4，避免 IPv6 超时问题
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	cutoffTime := time.Time{}
	if hoursAgo > 0 {
		cutoffTime = time.Now().Add(-time.Duration(hoursAgo) * time.Hour)
	}

	type ArticleWithMeta struct {
		Article   map[string]interface{}
		Source    string
		PubTime   time.Time
		SourceIdx int
	}

	var allArticles []ArticleWithMeta
	var sourceInfos []map[string]interface{}

	// 用于存储每个源的文章（分组）
	type SourceWithItems struct {
		FeedURL   string                   `json:"feed_url"`
		FeedTitle string                   `json:"feed_title"`
		Items     []map[string]interface{} `json:"items"`
		Error     string                   `json:"error,omitempty"`
	}
	var sourcesWithItems []SourceWithItems

	for idx, source := range sources {
		feed, err := fp.ParseURL(source.URL)
		if err != nil {
			// 某个源失败，记录但继续其他源
			sourceInfos = append(sourceInfos, map[string]interface{}{
				"feed_url":   source.URL,
				"feed_title": "解析失败",
				"item_count": 0,
				"error":      err.Error(),
			})
			sourcesWithItems = append(sourcesWithItems, SourceWithItems{
				FeedURL:   source.URL,
				FeedTitle: "解析失败",
				Items:     []map[string]interface{}{},
				Error:     err.Error(),
			})
			continue
		}

		// 解析关键词列表
		keywordList := []string{}
		if source.Keywords != "" {
			for _, kw := range strings.Split(source.Keywords, ",") {
				if trimmed := strings.TrimSpace(kw); trimmed != "" {
					keywordList = append(keywordList, strings.ToLower(trimmed))
				}
			}
		}

		itemCount := 0
		var sourceItems []map[string]interface{} // 当前源的文章列表
		for _, item := range feed.Items {
			// 时间过滤
			if !cutoffTime.IsZero() && item.PublishedParsed != nil {
				if item.PublishedParsed.Before(cutoffTime) {
					continue
				}
			}

			// 关键词过滤
			if len(keywordList) > 0 {
				matched := false
				searchText := strings.ToLower(item.Title + " " + item.Description)
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

			// 构建文章数据
			article := map[string]interface{}{
				"title":       item.Title,
				"link":        item.Link,
				"description": stripHTML(item.Description),
				"pub_date":    "",
				"author":      "",
				"categories":  []string{},
				"source":      feed.Title, // 标注来源
			}

			pubTime := time.Time{}
			if item.PublishedParsed != nil {
				article["pub_date"] = item.PublishedParsed.Format("2006-01-02 15:04:05")
				pubTime = *item.PublishedParsed
			} else if item.Published != "" {
				article["pub_date"] = item.Published
			}

			if item.Author != nil {
				article["author"] = item.Author.Name
			}

			if len(item.Categories) > 0 {
				article["categories"] = item.Categories
			}

			allArticles = append(allArticles, ArticleWithMeta{
				Article:   article,
				Source:    feed.Title,
				PubTime:   pubTime,
				SourceIdx: idx,
			})
			sourceItems = append(sourceItems, article) // 同时保存到当前源的列表
			itemCount++
		}

		sourceInfos = append(sourceInfos, map[string]interface{}{
			"feed_url":   source.URL,
			"feed_title": feed.Title,
			"item_count": itemCount,
		})

		// 保存当前源的分组数据
		sourcesWithItems = append(sourcesWithItems, SourceWithItems{
			FeedURL:   source.URL,
			FeedTitle: feed.Title,
			Items:     sourceItems,
		})
	}

	// 3. 去重
	dedupMap := make(map[string]bool)
	var dedupedArticles []ArticleWithMeta
	for _, am := range allArticles {
		key := ""
		if dedupBy == "link" {
			key = am.Article["link"].(string)
		} else {
			key = am.Article["title"].(string)
		}

		if !dedupMap[key] {
			dedupMap[key] = true
			dedupedArticles = append(dedupedArticles, am)
		}
	}

	// 4. 排序
	if sortBy == "time" {
		// 按发布时间倒序
		for i := 0; i < len(dedupedArticles); i++ {
			for j := i + 1; j < len(dedupedArticles); j++ {
				if dedupedArticles[i].PubTime.Before(dedupedArticles[j].PubTime) {
					dedupedArticles[i], dedupedArticles[j] = dedupedArticles[j], dedupedArticles[i]
				}
			}
		}
	}
	// sortBy == "source" 则保持原顺序

	// 5. 限制数量
	if len(dedupedArticles) > maxItems {
		dedupedArticles = dedupedArticles[:maxItems]
	}

	// 6. 提取最终文章列表
	var finalItems []map[string]interface{}
	for _, am := range dedupedArticles {
		finalItems = append(finalItems, am.Article)
	}

	// 7. 返回结果
	response := map[string]interface{}{
		"total_sources":      len(sources),
		"total_items":        len(finalItems),
		"sources":            sourceInfos,      // 统计信息
		"sources_with_items": sourcesWithItems, // 分组数据（按源分开）
		"items":              finalItems,       // 合并数据（所有文章）
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("从 %d 个订阅源获取了 %d 篇文章（去重后）", len(sources), len(finalItems)),
		Output: map[string]interface{}{
			"response":           response,
			"items":              finalItems,       // 合并数组（快捷访问）
			"sources_with_items": sourcesWithItems, // 分组数据（快捷访问）
			"total":              len(finalItems),
		},
	}, nil
}

// stripHTML 简单移除 HTML 标签
func stripHTML(s string) string {
	// 简单实现，移除 < > 之间的内容
	result := ""
	inTag := false
	for _, char := range s {
		if char == '<' {
			inTag = true
			continue
		}
		if char == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result += string(char)
		}
	}
	return strings.TrimSpace(result)
}

func float64Ptr(v float64) *float64 {
	return &v
}

func init() {
	utools.Register(NewRSSFeedTool())
}
