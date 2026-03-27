#!/bin/bash
# 即刻 MCP 服务器初始化脚本
set -e

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$SKILL_DIR"

export PATH="$HOME/go-sdk/go/bin:$HOME/go/bin:$PATH"

# 编译（如果没有二进制或源码更新了）
if [ ! -f "$SKILL_DIR/jike-mcp" ] || [ "$SKILL_DIR/main.go" -nt "$SKILL_DIR/jike-mcp" ]; then
    echo "编译 jike-mcp..."
    go build -o "$SKILL_DIR/jike-mcp" .
fi

# 启动（如果没在运行）
if ! curl -s http://localhost:8080/api/v1/status > /dev/null 2>&1; then
    echo "启动 jike-mcp 服务器..."
    "$SKILL_DIR/jike-mcp" --port 8080 --tokens "$SKILL_DIR/tokens.json" &
    sleep 2
fi

echo "即刻服务器就绪"
python3 "$SKILL_DIR/scripts/jike_client.py" status
