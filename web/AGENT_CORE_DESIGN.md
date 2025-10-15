# Agent 核心（agent_core）技术设计文档

> 目标：在现有工作流体系中，提供一个“受控、可审计、与模型无关”的智能代理（Agent）能力。Agent 能根据自然语言意图自主调用公司内部白名单工具（utools），并产出最终回答/结果，支持多轮对话与上下文记忆。

## 1. 模块概述
- 名称：`agent_core`
- 形态：作为新的 Tool 接入现有 utools 注册表与执行引擎（与 HTTP/JSON/Feishu 等工具同级）。
- 位置（建议）：`pkg/utools/agent/agent_tool.go`（后端），`web/src/components/tools/AgentConfig.vue`（前端面板）。
- 依赖：不引入大型外部 Agent 框架；仅使用模型 SDK、JSON 校验、Redis 作为可选组件。

## 2. 功能目标与非目标
### 2.1 功能目标
- 对话型任务：基于多轮上下文（结合 ContextManager）完成查询、汇总、通知等任务。
- 工具调用：在回合中选取白名单 utools 工具并执行，可多次调用，直至生成最终答案。
- 受控与审计：工具白名单、参数校验、超时/步数/预算限制；完整 trace 记录。
- 模型无关：采用统一的结构化文本协议（Action/Final JSON），不依赖 function calling。

### 2.2 非目标
- 不引入 LangChain/LlamaIndex 等大型框架；保持 Go 单栈与 utools 统一。
- 不追求通用市场化生态；以公司内部高价值用例为主，定制优先。

## 3. 架构与集成
```
┌───────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ ContextManager│ --> │   agent_core      │ --> │ ContextManager   │
│   (Prepare)   │     │  (Tool, ReAct循环)│     │    (Persist)     │
└───────────────┘     └──────────────────┘     └─────────────────┘
                                 │
                                 ▼
                        ┌──────────────────┐
                        │  utools Registry │ (HTTP/JSON/Feishu/Redis/...)
                        └──────────────────┘
```
- 引擎：保持不变；`agent_core` 作为普通 Tool 由引擎调度。
- 工具桥接（Tool Bridge）：`agent_core` 通过 utools 注册表获取白名单工具实例，调用其 `Execute`。
- 上下文：建议接入 `ContextManager(Prepare/Persist)`，在多轮对话中更稳。

## 4. 执行流程（ReAct + 结构化协议）
- 输入：
  - `messages_json`（推荐）来自 ContextManager；或
  - `prompt`（单轮模式）。
- 提示约束（System + Developer + Few-shot）：说明工具白名单、参数 schema、输出限制（仅 JSON）。
- 回合（最大步数 `max_steps`）：
  1) **模型输出** 一个 JSON：
     - Action：`{"type":"action","tool":"http_request","args":{...}}`
     - Final： `{"type":"final","answer":"..."}`
  2) **解析**：仅接受单个 JSON。解析失败 → 附带错误提示回灌 1 次“修复回合”。
  3) **执行**：按白名单与 schema 校验后调用 utools 工具，设置单步超时与限流；将 observation（摘要）回灌给模型。
  4) **终止**：遇 Final 或达 `max_steps`/`max_tool_calls`/预算上限。
- 输出：
  - `final_answer`（主结果）
  - `trace`（步骤、工具、参数摘要、observation 摘要、耗时、错误）
  - `used_tools`、`finish_reason`

## 5. 配置 Schema（建议）
```json
{
  "model": "gpt-4o-mini",
  "temperature": 0.7,
  "timeout": 120,
  "max_steps": 5,
  "max_tool_calls": 10,
  "allowed_tools": ["http_request","json_transform","redis_context","feishu_bot"],
  "single_call_timeout": 30,
  "total_token_budget": 0,
  "observation_max_len": 256,
  "primary_output": "final_answer"
}
```
- 必选：`model`、`max_steps`、`allowed_tools`。
- 安全：`single_call_timeout`、`max_tool_calls`、`total_token_budget`。
- 体验：`observation_max_len` 控制回灌文本长度，防上下文爆炸。

## 6. Action/Final 协议
- **Action**：`{"type":"action","tool":"<tool_code>","args":{...}}`
  - `tool` 必须在 `allowed_tools` 中。
  - `args` 依据对应工具的 `ConfigSchema` 校验（类型/必填/枚举/范围）。
- **Final**：`{"type":"final","answer":"..."}`
- 严格要求：模型“只输出一个 JSON”，不输出自然语言与 JSON 的混合；解析失败走一次修复回合。

## 7. Trace 结构（节点输出）
```json
{
  "final_answer": "...",
  "trace": [
    {
      "step": 1,
      "action": {"tool": "http_request", "args": {"url": "..."}},
      "observation": "200 OK, length=...",
      "elapsed_ms": 1200
    }
  ],
  "used_tools": {"http_request": {"count": 1, "total_ms": 1200}},
  "finish_reason": "final"
}
```

## 8. 与 ContextManager 的配合
- 推荐链路：`ContextManager(prepare)` → `agent_core(messages_json)` → `ContextManager(persist)`。
- ContextManager 已支持多厂商输出（openai/gemini/anthropic），Agent 默认读取中立 `messages_json` 即可。

## 9. 安全与治理
- 工具白名单：只开放安全工具；`http_request` 建议域名/方法白名单。
- 预算限制：`max_steps`、`max_tool_calls`、`single_call_timeout`、`total_token_budget`。
- 参数校验：严格使用工具 `ConfigSchema` 验证 `args`（必填/类型/范围）。
- 审计：`trace` 输出，必要时脱敏。
- 错误处理：解析失败/工具失败 → 回灌提示修复 1 次；仍失败则终止（防死循环）。

## 10. 前端面板（AgentConfig）
- 字段：模型、温度、最大步数、单步超时、白名单（多选）、预算限额、主输出字段（只读显示）。
- 辅助：变量助手（external/env/nodes），下拉选择白名单工具。
- 执行详情页：显示 `trace`（工具、参数摘要、observation、耗时）。

## 11. 测试方案
- 单测：
  - Action/Final 解析（含坏 JSON 与修复回合）。
  - 工具白名单与 `ConfigSchema` 参数校验。
  - 超时/限额终止逻辑；observation 摘要截断。
- 集成：
  - 单轮（仅 prompt）→ Final。
  - 多轮（配 ContextManager）→ 多步 Action → Final → Persist 追加历史。
  - 工具失败 → 修复回合 → Final；达上限 → finish_reason 正确。
- 多模型：OpenAI/Gemini/Claude 使用统一协议通过。

## 12. 监控与指标
- 指标：每步耗时、调用次数、finish_reason、失败率。
- 工具维度统计：count、avg_ms。
- Token/费用（可选）：从模型 SDK 中记录提示/完成 token 数。

## 13. 性能与限额
- observation 摘要长度控制；ContextManager 做历史窗口与 TTL。
- soft budget（token）与硬限制（步数/次数/时长）。

## 14. 上线与回滚
- MVP：OpenAI + 3–5 个白名单工具 → 业务试点。
- Trace 观察与 few-shot 优化。
- 回滚：`agent_core` 为独立 Tool；下线节点或恢复旧工作流即可。

## 15. 提示词模板（建议）
### 15.1 System（约束）
- 你可以调用以下工具（列出 tool_code + 关键参数说明）。
- 只能输出一个 JSON，不要输出自然语言与 JSON 的混合内容。
- 两种输出之一：
  - Action: {"type":"action","tool":"<whitelist_tool>","args":{...}}
  - Final:  {"type":"final","answer":"..."}
- 如果工具执行失败或参数无效，你会收到错误提示；请纠正并重试或返回 Final。

### 15.2 Developer（上下文）
- 任务目标、限制（域名白名单、时间限制）。
- 仅提供 observation 摘要，若信息不足请再次合理调用工具获取详情。

### 15.3 Few-shot 示例（至少 2–3 组）
- 成功 action → final；工具失败 → 修复后 action → final；无工具直接 final。

## 16. 开放问题（决策后可实施）
- 是否提供“调用子工作流”的受控函数（`invoke_workflow`）以复用复杂任务链？
- 是否在 UI 中增加“provider 输出勾选开关”（仅输出所需格式减少变量噪音）？
- 是否把 `trace` 长期入库用于审计？

## 17. 验收标准
- `agent_core` 可在 `messages_json`/`prompt` 两种模式工作。
- Action/Final 协议稳定，解析与修复回合可靠，错误时不崩溃。
- 工具白名单 + 参数校验生效；`http_request` 受域名/方法限制（如配置）。
- 与 ContextManager 配合，能正确读取/写回会话历史（JSON 数组）。
- 执行详情可直观看到 `trace` 与 `final_answer`。
- 外部同步调用可直接读取 `outputs.agent_core.final_answer`。

---

> 备注：本设计优先保障“公司内部可控与维护成本低”。若未来需要引入框架（LangChain 等），建议以“边车服务”形式集成，由 `agent_core` 统一管理工具调用与审计，避免工程复杂度失控。

