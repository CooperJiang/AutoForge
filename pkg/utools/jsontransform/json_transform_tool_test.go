package jsontransform

import (
	"auto-forge/pkg/utools"
	"context"
	"testing"
)

func TestJSONTransform_JSExpression(t *testing.T) {
	tool := NewJSONTransformTool()
	config := map[string]interface{}{
		"data_source": `[{"url":"https://a.com"},{"url":"https://b.com"}]`,
		"expression":  "data.map(item => item.url)",
		"output_name": "urls",
		"timeout_ms":  1500.0,
	}

	ctx := &utools.ExecutionContext{Context: context.Background()}

	result, err := tool.Execute(ctx, config)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got: %s", result.Message)
	}

	outputs := result.Output["urls"].([]interface{})
	if len(outputs) != 2 {
		t.Fatalf("expected 2 urls, got %d", len(outputs))
	}
	if outputs[0].(string) != "https://a.com" {
		t.Fatalf("unexpected first url: %v", outputs[0])
	}
}

func TestJSONTransform_ContextLookup(t *testing.T) {
	tool := NewJSONTransformTool()
	config := map[string]interface{}{
		"data_source": `{{nodes.prev.output}}`,
		"expression":  "data.filter(item => item.ok).map(item => item.value)",
	}

	ctx := &utools.ExecutionContext{
		Context: context.Background(),
		Variables: map[string]interface{}{
			"nodes": map[string]interface{}{
				"prev": map[string]interface{}{
					"output": []interface{}{
						map[string]interface{}{"ok": true, "value": "a"},
						map[string]interface{}{"ok": false, "value": "b"},
					},
				},
			},
		},
	}

	result, err := tool.Execute(ctx, config)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got: %s", result.Message)
	}

	values := result.Output["result"].([]interface{})
	if len(values) != 1 || values[0].(string) != "a" {
		t.Fatalf("unexpected result: %v", values)
	}
}
