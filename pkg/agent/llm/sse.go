package llm

import (
	"bufio"
	"io"
	"strings"
)

// SSEEvent SSE 事件
type SSEEvent struct {
	Event string
	Data  string
	ID    string
}

// SSEReader SSE 读取器
type SSEReader struct {
	scanner *bufio.Scanner
}

// NewSSEReader 创建 SSE 读取器
func NewSSEReader(r io.Reader) *SSEReader {
	return &SSEReader{
		scanner: bufio.NewScanner(r),
	}
}

// Read 读取下一个事件
func (r *SSEReader) Read() (*SSEEvent, error) {
	event := &SSEEvent{}

	for r.scanner.Scan() {
		line := r.scanner.Text()

		// 空行表示事件结束
		if line == "" {
			if event.Data != "" {
				return event, nil
			}
			continue
		}

		// 解析字段
		if strings.HasPrefix(line, "event:") {
			event.Event = strings.TrimSpace(line[6:])
		} else if strings.HasPrefix(line, "data:") {
			if event.Data != "" {
				event.Data += "\n"
			}
			event.Data += strings.TrimSpace(line[5:])
		} else if strings.HasPrefix(line, "id:") {
			event.ID = strings.TrimSpace(line[3:])
		}
		// 忽略注释行（以 : 开头）
	}

	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}
