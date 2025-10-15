#!/bin/bash

# 测试工作流输出格式脚本
# 用法: ./test_workflow_output.sh <workflow_id>

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:7777"
WORKFLOW_ID="${1}"
MAX_RETRIES=30
RETRY_INTERVAL=2

if [ -z "$WORKFLOW_ID" ]; then
    echo -e "${RED}错误: 请提供工作流 ID${NC}"
    echo "用法: $0 <workflow_id>"
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}测试工作流输出格式${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 测试 1: 首次对话
echo -e "${YELLOW}[测试 1] 首次对话 - 介绍自己${NC}"
echo "发送消息: '你好，我叫张三，我是一名程序员'"
echo ""

RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/workflows/${WORKFLOW_ID}/execute" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -d '{
    "params": {
      "session_id": "test_user_001",
      "user_message": "你好，我叫张三，我是一名程序员"
    }
  }')

EXECUTION_ID=$(echo "$RESPONSE" | jq -r '.data.execution_id')
echo -e "${GREEN}✓ 执行ID: $EXECUTION_ID${NC}"
echo ""

# 等待执行完成
echo -e "${YELLOW}等待执行完成...${NC}"
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    sleep $RETRY_INTERVAL

    EXEC_RESULT=$(curl -s -X GET "${BASE_URL}/api/v1/executions/${EXECUTION_ID}" \
      -H "Authorization: Bearer demo-token")

    STATUS=$(echo "$EXEC_RESULT" | jq -r '.data.status')

    if [ "$STATUS" = "success" ]; then
        echo -e "${GREEN}✓ 执行成功！${NC}"
        break
    elif [ "$STATUS" = "failed" ]; then
        echo -e "${RED}✗ 执行失败${NC}"
        echo "$EXEC_RESULT" | jq '.data.error'
        exit 1
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -ne "${YELLOW}等待中... ($RETRY_COUNT/$MAX_RETRIES)\r${NC}"
done

if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo -e "${RED}✗ 执行超时${NC}"
    exit 1
fi

echo ""

# 提取最后一个节点的输出
FINAL_OUTPUT=$(echo "$EXEC_RESULT" | jq '.data.node_logs[-1].output')

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}最终输出 (格式化后)${NC}"
echo -e "${BLUE}========================================${NC}"
echo "$FINAL_OUTPUT" | jq '.'
echo ""

# 提取关键字段
AI_REPLY=$(echo "$FINAL_OUTPUT" | jq -r '.data.reply')
MESSAGE_COUNT=$(echo "$FINAL_OUTPUT" | jq -r '.data.message_count')

echo -e "${GREEN}AI 回复:${NC} $AI_REPLY"
echo -e "${GREEN}消息数量:${NC} $MESSAGE_COUNT"
echo ""

# 测试 2: 记忆测试
echo -e "${BLUE}========================================${NC}"
echo -e "${YELLOW}[测试 2] 记忆测试 - 询问姓名和职业${NC}"
echo "发送消息: '我叫什么名字？我的职业是什么？'"
echo ""

RESPONSE2=$(curl -s -X POST "${BASE_URL}/api/v1/workflows/${WORKFLOW_ID}/execute" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -d '{
    "params": {
      "session_id": "test_user_001",
      "user_message": "我叫什么名字？我的职业是什么？"
    }
  }')

EXECUTION_ID2=$(echo "$RESPONSE2" | jq -r '.data.execution_id')
echo -e "${GREEN}✓ 执行ID: $EXECUTION_ID2${NC}"
echo ""

# 等待执行完成
echo -e "${YELLOW}等待执行完成...${NC}"
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    sleep $RETRY_INTERVAL

    EXEC_RESULT2=$(curl -s -X GET "${BASE_URL}/api/v1/executions/${EXECUTION_ID2}" \
      -H "Authorization: Bearer demo-token")

    STATUS2=$(echo "$EXEC_RESULT2" | jq -r '.data.status')

    if [ "$STATUS2" = "success" ]; then
        echo -e "${GREEN}✓ 执行成功！${NC}"
        break
    elif [ "$STATUS2" = "failed" ]; then
        echo -e "${RED}✗ 执行失败${NC}"
        echo "$EXEC_RESULT2" | jq '.data.error'
        exit 1
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -ne "${YELLOW}等待中... ($RETRY_COUNT/$MAX_RETRIES)\r${NC}"
done

echo ""

# 提取最后一个节点的输出
FINAL_OUTPUT2=$(echo "$EXEC_RESULT2" | jq '.data.node_logs[-1].output')

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}最终输出 (格式化后)${NC}"
echo -e "${BLUE}========================================${NC}"
echo "$FINAL_OUTPUT2" | jq '.'
echo ""

# 提取关键字段
AI_REPLY2=$(echo "$FINAL_OUTPUT2" | jq -r '.data.reply')
MESSAGE_COUNT2=$(echo "$FINAL_OUTPUT2" | jq -r '.data.message_count')

echo -e "${GREEN}AI 回复:${NC} $AI_REPLY2"
echo -e "${GREEN}消息数量:${NC} $MESSAGE_COUNT2"
echo ""

# 验证记忆功能
if [[ "$AI_REPLY2" == *"张三"* ]] && [[ "$AI_REPLY2" == *"程序员"* ]]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✓ 记忆功能测试通过！${NC}"
    echo -e "${GREEN}AI 成功记住了姓名和职业${NC}"
    echo -e "${GREEN}========================================${NC}"
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}✗ 记忆功能测试失败${NC}"
    echo -e "${RED}AI 没有正确记住信息${NC}"
    echo -e "${RED}========================================${NC}"
fi

echo ""

# 测试 3: 会话隔离测试
echo -e "${BLUE}========================================${NC}"
echo -e "${YELLOW}[测试 3] 会话隔离 - 新用户${NC}"
echo "发送消息: '我叫什么名字？'"
echo ""

RESPONSE3=$(curl -s -X POST "${BASE_URL}/api/v1/workflows/${WORKFLOW_ID}/execute" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -d '{
    "params": {
      "session_id": "test_user_002",
      "user_message": "我叫什么名字？"
    }
  }')

EXECUTION_ID3=$(echo "$RESPONSE3" | jq -r '.data.execution_id')
echo -e "${GREEN}✓ 执行ID: $EXECUTION_ID3${NC}"
echo ""

# 等待执行完成
echo -e "${YELLOW}等待执行完成...${NC}"
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    sleep $RETRY_INTERVAL

    EXEC_RESULT3=$(curl -s -X GET "${BASE_URL}/api/v1/executions/${EXECUTION_ID3}" \
      -H "Authorization: Bearer demo-token")

    STATUS3=$(echo "$EXEC_RESULT3" | jq -r '.data.status')

    if [ "$STATUS3" = "success" ]; then
        echo -e "${GREEN}✓ 执行成功！${NC}"
        break
    elif [ "$STATUS3" = "failed" ]; then
        echo -e "${RED}✗ 执行失败${NC}"
        echo "$EXEC_RESULT3" | jq '.data.error'
        exit 1
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -ne "${YELLOW}等待中... ($RETRY_COUNT/$MAX_RETRIES)\r${NC}"
done

echo ""

# 提取最后一个节点的输出
FINAL_OUTPUT3=$(echo "$EXEC_RESULT3" | jq '.data.node_logs[-1].output')

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}最终输出 (格式化后)${NC}"
echo -e "${BLUE}========================================${NC}"
echo "$FINAL_OUTPUT3" | jq '.'
echo ""

AI_REPLY3=$(echo "$FINAL_OUTPUT3" | jq -r '.data.reply')
echo -e "${GREEN}AI 回复:${NC} $AI_REPLY3"
echo ""

# 验证会话隔离
if [[ "$AI_REPLY3" != *"张三"* ]]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✓ 会话隔离测试通过！${NC}"
    echo -e "${GREEN}新用户的对话不包含之前的信息${NC}"
    echo -e "${GREEN}========================================${NC}"
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}✗ 会话隔离测试失败${NC}"
    echo -e "${RED}新用户的对话泄露了其他用户的信息${NC}"
    echo -e "${RED}========================================${NC}"
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}测试完成！${NC}"
echo -e "${BLUE}========================================${NC}"
