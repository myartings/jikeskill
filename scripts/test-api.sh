#!/bin/bash
# 即刻 MCP 服务器 API 测试脚本
# 用法: ./scripts/test-api.sh [port]

PORT=${1:-9999}
BASE="http://localhost:$PORT"
PASS=0
FAIL=0

test_api() {
  local name="$1"
  local method="$2"
  local path="$3"
  local data="$4"
  local expect="$5"  # 期望在响应中包含的字符串

  if [ "$method" = "GET" ]; then
    resp=$(curl -s "$BASE$path")
  else
    resp=$(curl -s -X POST "$BASE$path" -H "Content-Type: application/json" -d "$data")
  fi

  if echo "$resp" | grep -q "$expect"; then
    echo "  ✓ $name"
    PASS=$((PASS + 1))
  else
    echo "  ✗ $name"
    echo "    响应: $(echo "$resp" | head -c 200)"
    FAIL=$((FAIL + 1))
  fi
}

echo "测试 Jike MCP Server (端口 $PORT)"
echo "================================"

# 检查服务器是否运行
if ! curl -s "$BASE/api/v1/status" > /dev/null 2>&1; then
  echo "✗ 服务器未运行在端口 $PORT"
  echo "  请先启动: ./jike-mcp --port $PORT"
  exit 1
fi

echo ""
echo "[认证]"
test_api "登录状态" GET "/api/v1/status" "" "logged_in"

# 检查是否已登录
STATUS=$(curl -s "$BASE/api/v1/status")
if echo "$STATUS" | grep -q '"logged_in":true'; then
  LOGGED_IN=true
  echo "  (已登录，测试全部接口)"
else
  LOGGED_IN=false
  echo "  (未登录，跳过需认证的接口)"
fi

echo ""
echo "[登录流程]"
test_api "获取二维码" POST "/api/v1/login/qrcode" "{}" "uuid"

if [ "$LOGGED_IN" = true ]; then
  echo ""
  echo "[动态流]"
  test_api "关注动态" POST "/api/v1/feeds/following" "{}" "data"
  test_api "推荐动态" POST "/api/v1/feeds/recommend" "{}" "data"

  echo ""
  echo "[搜索]"
  test_api "搜索" POST "/api/v1/search" '{"keyword":"AI"}' "data"

  echo ""
  echo "[用户]"
  # 用当前登录用户的 username 测试
  USERNAME=$(echo "$STATUS" | python3 -c "import sys,json; print(json.load(sys.stdin).get('user',{}).get('username',''))" 2>/dev/null)
  if [ -n "$USERNAME" ]; then
    test_api "用户资料" GET "/api/v1/user/$USERNAME" "" "screenName"
    test_api "用户帖子" POST "/api/v1/user/$USERNAME/posts" "{}" "data"
  fi
fi

echo ""
echo "================================"
echo "结果: $PASS 通过, $FAIL 失败"

if [ $FAIL -gt 0 ]; then
  exit 1
fi
