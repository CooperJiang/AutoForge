package prompt

import (
	"fmt"
	"strings"
)

// Template 提示词模板
type Template struct {
	Name     string
	Template string
}

// Render 渲染模板
func (t *Template) Render(vars map[string]string) string {
	result := t.Template
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

// 系统提示词模板
var SystemPrompt = Template{
	Name: "system",
	Template: `You are a helpful AI assistant with access to various tools. Your goal is to help users accomplish their tasks by intelligently using the available tools.

Key principles:
1. Think step by step before taking action
2. Use tools when necessary to gather information or perform actions
3. Provide clear and helpful responses based on tool results
4. If you're unsure, ask for clarification
5. Always explain your reasoning

Remember: You can call multiple tools in sequence to accomplish complex tasks.`,
}

// ReAct 提示词模板
var ReActPrompt = Template{
	Name: "react",
	Template: `You have access to the following tools:

{tool_definitions}

To use a tool, respond with a function call in the format provided by the API.

Previous conversation:
{conversation_history}

User request: {user_message}

Think carefully about what tools you need to use and in what order. Let's work through this step by step.`,
}

// Plan 提示词模板
var PlanPrompt = Template{
	Name: "plan",
	Template: `You are a planning AI. Create a detailed execution plan for the user's request.

Available tools (YOU MUST USE THESE EXACT NAMES):
{tool_definitions}

User request: {user_message}

CRITICAL RULES:
1. You MUST ONLY use tool names from the list above - DO NOT invent or modify tool names
2. Each tool name must EXACTLY match one from the available tools list
3. If you use a tool name not in the list, the execution will FAIL
4. When the user asks for external data, you MUST use data-fetching tools FIRST
5. DO NOT make up or fabricate data - always fetch real data using tools
6. Break down complex tasks into multiple steps with the right tools
7. Read each tool's description and parameters carefully to understand what it does
8. Pay attention to parameter types (string, object, boolean, etc.) and required vs optional
9. If a tool accepts multiple input types, choose the most efficient one based on previous step outputs

PLANNING GUIDELINES:
- Analyze the available tools and their capabilities from the descriptions above
- Choose tools based on their descriptions, parameters, and output schemas
- Consider the output of each step when planning the next step
- Use the most direct path to accomplish the task
- If a step depends on data from a previous step, make sure that data will be available in the output

VERIFICATION CHECKLIST:
✓ Every "tool" field uses an EXACT name from the available tools list
✓ No invented or modified tool names
✓ If task requires external data, include data-fetching step FIRST
✓ Steps are in logical order with proper data flow
✓ Each step has clear description and reasoning
✓ Parameter types will match the tool's schema
✓ Return ONLY valid JSON, no explanations
✓ The response has a "steps" array

Now create the plan (remember: use EXACT tool names and fetch real data when needed):`,
}

// Summary 提示词模板
var SummaryPrompt = Template{
	Name: "summary",
	Template: `Based on the following tool execution results, provide a clear and helpful final answer to the user's request.

User request: {user_message}

Execution steps and results:
{execution_trace}

Provide a natural language summary that:
1. Directly answers the user's question
2. Highlights key information from the tool results
3. If any step failed, clearly explain what went wrong and why
4. If the task couldn't be completed, explain what was attempted and suggest alternatives
5. Is clear, honest, and helpful
6. Avoids unnecessary technical details

IMPORTANT: If any tool execution failed, you MUST:
- Acknowledge the failure clearly
- Explain which step failed and why
- Tell the user what information is missing or unavailable
- Suggest what they can do instead

Final answer:`,
}

// Error 提示词模板
var ErrorPrompt = Template{
	Name: "error",
	Template: `An error occurred while executing a tool:

Tool: {tool_name}
Error: {error_message}

Previous context:
{context}

Suggest an alternative approach or explain why the task cannot be completed. Be helpful and constructive.`,
}

// FormatToolDefinitions 格式化工具定义为文本
func FormatToolDefinitions(tools []map[string]interface{}) string {
	var builder strings.Builder

	builder.WriteString("=== AVAILABLE TOOLS (USE EXACT NAMES) ===\n\n")

	for i, tool := range tools {
		if function, ok := tool["function"].(map[string]interface{}); ok {
			builder.WriteString(formatSingleTool(function, i+1))
			if i < len(tools)-1 {
				builder.WriteString("\n")
				builder.WriteString("---\n\n")
			}
		}
	}

	return builder.String()
}

// formatSingleTool 格式化单个工具定义
func formatSingleTool(function map[string]interface{}, index int) string {
	var builder strings.Builder

	// 序号和工具名称
	if name, ok := function["name"].(string); ok {
		builder.WriteString(fmt.Sprintf("%d. Tool Name: %s\n", index, name))
	}

	// 描述
	if desc, ok := function["description"].(string); ok {
		builder.WriteString("   Description: ")
		builder.WriteString(desc)
		builder.WriteString("\n")
	}

	// 参数
	if params, ok := function["parameters"].(map[string]interface{}); ok {
		// 必需参数
		var required []string
		if req, ok := params["required"].([]interface{}); ok {
			for _, r := range req {
				if str, ok := r.(string); ok {
					required = append(required, str)
				}
			}
		} else if req, ok := params["required"].([]string); ok {
			required = req
		}

		// 参数列表
		if props, ok := params["properties"].(map[string]interface{}); ok {
			builder.WriteString("   Parameters:\n")
			for key, value := range props {
				if propMap, ok := value.(map[string]interface{}); ok {
					// 检查是否必需
					isRequired := false
					for _, r := range required {
						if r == key {
							isRequired = true
							break
						}
					}

					builder.WriteString("     - ")
					builder.WriteString(key)

					// 类型
					if propType, ok := propMap["type"].(string); ok {
						builder.WriteString(" (")
						builder.WriteString(propType)
						if isRequired {
							builder.WriteString(", required")
						} else {
							builder.WriteString(", optional")
						}
						builder.WriteString(")")
					}

					// 描述
					if desc, ok := propMap["description"].(string); ok {
						builder.WriteString(": ")
						builder.WriteString(desc)
					}

					// 枚举值
					if enum, ok := propMap["enum"].([]interface{}); ok {
						builder.WriteString(" [")
						for i, e := range enum {
							if i > 0 {
								builder.WriteString(", ")
							}
							builder.WriteString(fmt.Sprintf("%v", e))
						}
						builder.WriteString("]")
					}

					// 默认值
					if def, ok := propMap["default"]; ok {
						builder.WriteString(fmt.Sprintf(" (default: %v)", def))
					}

					builder.WriteString("\n")
				}
			}
		}
	}

	// 输出字段（从 metadata 中获取）
	if metadata, ok := function["metadata"].(map[string]interface{}); ok {
		if outputSchema, ok := metadata["output_fields_schema"].(map[string]interface{}); ok && len(outputSchema) > 0 {
			builder.WriteString("   Output Fields (how to extract from previous step):\n")

			// 按字段名排序，确保一致性
			var keys []string
			for key := range outputSchema {
				keys = append(keys, key)
			}

			for _, key := range keys {
				field := outputSchema[key]
				if fieldDef, ok := field.(map[string]interface{}); ok {
					label := ""
					if l, ok := fieldDef["label"].(string); ok {
						label = l
					}

					fieldType := ""
					if t, ok := fieldDef["type"].(string); ok {
						fieldType = t
					}

					builder.WriteString(fmt.Sprintf("     - %s (%s): %s\n", key, fieldType, label))
					builder.WriteString(fmt.Sprintf("       Access: output.%s\n", key))
				}
			}
		}
	}

	return builder.String()
}

// FormatConversationHistory 格式化对话历史
func FormatConversationHistory(history string) string {
	if history == "" {
		return "(No previous conversation)"
	}
	return history
}

// FormatExecutionTrace 格式化执行轨迹
func FormatExecutionTrace(steps []map[string]interface{}) string {
	var builder strings.Builder

	for i, step := range steps {
		builder.WriteString("Step ")
		builder.WriteString(string(rune('0' + i + 1)))
		builder.WriteString(":\n")

		if action, ok := step["action"].(map[string]interface{}); ok {
			if tool, ok := action["tool"].(string); ok {
				builder.WriteString("  Tool: ")
				builder.WriteString(tool)
				builder.WriteString("\n")
			}
		}

		if obs, ok := step["observation"].(string); ok {
			builder.WriteString("  Result: ")
			builder.WriteString(obs)
			builder.WriteString("\n")
		}

		if i < len(steps)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
