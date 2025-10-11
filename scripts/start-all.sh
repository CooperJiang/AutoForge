#!/bin/bash
echo "🚀 启动开发环境"
echo ""
echo "正在启动服务..."
echo ""

# 使用Makefile传入的项目根目录
PROJECT_ROOT="/Users/lilithgames/Desktop/cooper-core/gin-template"

# 启动后端
echo "🔧 启动后端 (Go) - 端口 7777"
osascript -e "tell application \"Terminal\" to do script \"cd \\"$PROJECT_ROOT\\" && ./scripts/start-backend.sh\"" >/dev/null 2>&1 &
sleep 1

# 启动用户端
echo "🌐 启动用户端 (Vue3) - 端口 3200"
osascript -e "tell application \"Terminal\" to do script \"cd \\"$PROJECT_ROOT\\" && ./scripts/start-web.sh\"" >/dev/null 2>&1 &
sleep 1

echo ""
echo "✅ 开发环境已在独立终端中启动！"
echo ""
echo "📱 服务访问地址:"
echo "  🎯 后端API:  http://localhost:7777"
echo "  🌐 用户端:   http://localhost:3200"
echo ""
echo "💡 提示: 各服务在独立的终端窗口中运行"
echo "💡 关闭: 在各终端窗口中按 Ctrl+C 停止服务"

