#!/bin/bash
# 即刻 MCP 服务器初始化脚本
set -e

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PORT=8181

# 启动（如果没在运行）
if ! curl -s "http://localhost:$PORT/api/v1/status" > /dev/null 2>&1; then
    echo "启动 jike-mcp 服务器..."
    nohup "$SKILL_DIR/jike-mcp" --port "$PORT" --tokens "$SKILL_DIR/tokens.json" \
        > /tmp/jike-mcp.log 2>&1 &
    sleep 2
fi

echo "即刻服务器就绪 (port $PORT)"
JIKE_API_URL="http://localhost:$PORT" python3 "$SKILL_DIR/scripts/jike_client.py" status
