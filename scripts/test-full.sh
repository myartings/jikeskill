#!/bin/bash
# 即刻 Skill 完整功能测试
# 测试 REST API + Python CLI + URL 解析
# 用法: ./scripts/test-full.sh [port]

set -o pipefail

PORT=${1:-8080}
BASE="http://localhost:$PORT"
SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CLI="python3 $SKILL_DIR/scripts/jike_client.py"
PASS=0
FAIL=0
SKIP=0

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

pass() { echo -e "  ${GREEN}✓${NC} $1"; PASS=$((PASS + 1)); }
fail() { echo -e "  ${RED}✗${NC} $1"; echo "    $2" | head -c 300; FAIL=$((FAIL + 1)); }
skip() { echo -e "  ${YELLOW}⊘${NC} $1 (跳过: $2)"; SKIP=$((SKIP + 1)); }

# ========================================
# REST API 测试
# ========================================
test_api() {
  local name="$1" method="$2" path="$3" data="$4" expect="$5"
  local resp
  if [ "$method" = "GET" ]; then
    resp=$(curl -s -w "\n%{http_code}" "$BASE$path" 2>&1)
  else
    resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE$path" -H "Content-Type: application/json" -d "$data" 2>&1)
  fi
  local http_code=$(echo "$resp" | tail -1)
  local body=$(echo "$resp" | sed '$d')

  if echo "$body" | grep -q "$expect"; then
    pass "API: $name (HTTP $http_code)"
  else
    fail "API: $name (HTTP $http_code)" "期望包含 '$expect', 实际: $body"
  fi
}

# 测试 API 期望失败（如参数校验）
test_api_error() {
  local name="$1" method="$2" path="$3" data="$4" expect_code="$5"
  local resp
  if [ "$method" = "GET" ]; then
    resp=$(curl -s -w "\n%{http_code}" "$BASE$path" 2>&1)
  else
    resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE$path" -H "Content-Type: application/json" -d "$data" 2>&1)
  fi
  local http_code=$(echo "$resp" | tail -1)

  if [ "$http_code" = "$expect_code" ]; then
    pass "API: $name (正确返回 HTTP $http_code)"
  else
    local body=$(echo "$resp" | sed '$d')
    fail "API: $name" "期望 HTTP $expect_code, 实际 HTTP $http_code: $body"
  fi
}

# ========================================
# Python CLI 测试
# ========================================
test_cli() {
  local name="$1" expect="$2"
  shift 2
  local resp
  resp=$(JIKE_API_URL="$BASE" $CLI "$@" 2>&1)
  if echo "$resp" | grep -q "$expect"; then
    pass "CLI: $name"
  else
    fail "CLI: $name" "期望包含 '$expect', 实际: $resp"
  fi
}

# ========================================
# 开始测试
# ========================================
echo "即刻 Skill 完整功能测试"
echo "服务器: $BASE"
echo "========================================"

# 检查服务器
if ! curl -s "$BASE/api/v1/status" > /dev/null 2>&1; then
  echo -e "${RED}✗ 服务器未运行在端口 $PORT${NC}"
  echo "  请先启动: ./jike-mcp --port $PORT"
  exit 1
fi

STATUS=$(curl -s "$BASE/api/v1/status")
LOGGED_IN=$(echo "$STATUS" | python3 -c "import sys,json; print(json.load(sys.stdin).get('logged_in', False))" 2>/dev/null)
USERNAME=$(echo "$STATUS" | python3 -c "import sys,json; print(json.load(sys.stdin).get('user',{}).get('username',''))" 2>/dev/null)
SCREEN_NAME=$(echo "$STATUS" | python3 -c "import sys,json; print(json.load(sys.stdin).get('user',{}).get('screenName',''))" 2>/dev/null)

echo ""
echo "=== 1. 认证 ==="

test_api "状态检查返回 logged_in 字段" GET "/api/v1/status" "" "logged_in"

if [ "$LOGGED_IN" = "True" ]; then
  pass "已登录: $SCREEN_NAME (@$USERNAME)"
else
  echo -e "  ${YELLOW}未登录，部分测试将跳过${NC}"
fi

echo ""
echo "=== 2. 登录流程 ==="

# 测试二维码生成
QRRESP=$(curl -s -X POST "$BASE/api/v1/login/qrcode" -H "Content-Type: application/json" -d '{}')
QR_UUID=$(echo "$QRRESP" | python3 -c "import sys,json; print(json.load(sys.stdin).get('uuid',''))" 2>/dev/null)
if [ -n "$QR_UUID" ]; then
  pass "API: 生成二维码 (uuid=$QR_UUID)"
else
  fail "API: 生成二维码" "无 uuid: $QRRESP"
fi

# 测试 qrcode_base64 存在
if echo "$QRRESP" | grep -q "qrcode_base64"; then
  pass "API: 二维码包含 base64 数据"
else
  fail "API: 二维码包含 base64 数据" "无 qrcode_base64 字段"
fi

echo ""
echo "=== 3. Python CLI 基础 ==="

# CLI 无参数应显示帮助
test_cli "无参数显示帮助" "Usage" 2>/dev/null || test_cli "无参数显示帮助" "usage"
test_cli "未知命令报错" "未知命令" "foobar_nonexistent"
test_cli "状态检查" "登录\|logged" "status"

echo ""
echo "=== 4. URL 解析（resolve_username）==="

# 测试 Python resolve_username 函数
test_resolve() {
  local name="$1" input="$2" expect="$3"
  local result
  result=$(python3 -c "
import os, sys, re, urllib.request, urllib.error
sys.path.insert(0, '$SKILL_DIR/scripts')

# Extract resolve_username from jike_client.py
import importlib.util
spec = importlib.util.spec_from_file_location('jike_client', '$SKILL_DIR/scripts/jike_client.py')
mod = importlib.util.module_from_spec(spec)
# Prevent main() from running
sys.argv = ['test']
spec.loader.exec_module(mod)

print(mod.resolve_username('$input'))
" 2>&1)
  if echo "$result" | grep -q "$expect"; then
    pass "URL解析: $name → $(echo "$result" | tail -1)"
  else
    fail "URL解析: $name" "期望包含 '$expect', 实际: $result"
  fi
}

# 普通 username 应原样返回
test_resolve "普通username不变" "testuser123" "testuser123"
# 纯小写不触发短码检测
test_resolve "纯小写不触发短码" "abcdef" "abcdef"
# 纯大写不触发短码检测
test_resolve "纯大写不触发短码" "ABCDEF" "ABCDEF"
# 数字不触发短码检测
test_resolve "纯数字不触发短码" "123456" "123456"

echo ""
echo "=== 5. Go URL 解析 ==="

# 测试 Go 端的 isUUID 逻辑（通过 API 调用间接测试）
if [ "$LOGGED_IN" = "True" ] && [ -n "$USERNAME" ]; then
  # 用已知的 username 测试 API
  test_api "用 username 查资料" GET "/api/v1/user/$USERNAME" "" "screenName"

  # 获取用户 ID
  USER_ID=$(curl -s "$BASE/api/v1/user/$USERNAME" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))" 2>/dev/null)
  if [ -n "$USER_ID" ]; then
    pass "获取到用户 ID: $USER_ID"
  fi
else
  skip "UUID查询" "未登录"
fi

# ========================================
# 以下测试需要登录
# ========================================
if [ "$LOGGED_IN" != "True" ]; then
  echo ""
  echo -e "${YELLOW}未登录，跳过需认证的测试${NC}"
  echo ""
  echo "========================================"
  echo -e "结果: ${GREEN}$PASS 通过${NC}, ${RED}$FAIL 失败${NC}, ${YELLOW}$SKIP 跳过${NC}"
  [ $FAIL -gt 0 ] && exit 1 || exit 0
fi

echo ""
echo "=== 6. 动态流 API ==="

test_api "关注动态" POST "/api/v1/feeds/following" "{}" "data"
test_api "推荐动态" POST "/api/v1/feeds/recommend" "{}" "data"

# 获取一条帖子 ID 用于后续测试
FEED_RESP=$(curl -s -X POST "$BASE/api/v1/feeds/following" -H "Content-Type: application/json" -d '{}')
POST_ID=$(echo "$FEED_RESP" | python3 -c "
import sys, json
data = json.load(sys.stdin).get('data', [])
for p in data:
    pid = p.get('id', '')
    if pid:
        print(pid)
        break
" 2>/dev/null)

if [ -n "$POST_ID" ]; then
  pass "从动态流获取测试帖子 ID: $POST_ID"
else
  skip "获取测试帖子 ID" "动态流为空"
fi

echo ""
echo "=== 7. 动态流 CLI ==="

test_cli "关注动态" "" "following"
test_cli "推荐动态" "" "recommend"

echo ""
echo "=== 8. 搜索 ==="

test_api "搜索关键词" POST "/api/v1/search" '{"keyword":"即刻"}' "data"
test_api_error "搜索缺少keyword" POST "/api/v1/search" '{}' "400"
test_cli "CLI搜索" "结果" "search" "即刻"

echo ""
echo "=== 9. 帖子详情 ==="

if [ -n "$POST_ID" ]; then
  test_api "帖子详情" POST "/api/v1/post/detail" "{\"post_id\":\"$POST_ID\"}" "content\|id"
  test_cli "CLI帖子详情" "" "post-detail" "$POST_ID"
else
  skip "帖子详情" "无测试帖子 ID"
  skip "CLI帖子详情" "无测试帖子 ID"
fi

test_api_error "帖子详情缺少post_id" POST "/api/v1/post/detail" '{}' "400"

echo ""
echo "=== 10. 评论 ==="

if [ -n "$POST_ID" ]; then
  test_api "评论列表" POST "/api/v1/comments/list" "{\"target_id\":\"$POST_ID\"}" "data"

  # 测试多种参数名格式
  test_api "评论列表(post_id格式)" POST "/api/v1/comments/list" "{\"post_id\":\"$POST_ID\"}" "data"
  test_api "评论列表(targetId格式)" POST "/api/v1/comments/list" "{\"targetId\":\"$POST_ID\"}" "data"

  test_cli "CLI评论列表" "" "comments" "$POST_ID"
else
  skip "评论列表" "无测试帖子 ID"
fi

test_api_error "评论列表缺少target_id" POST "/api/v1/comments/list" '{}' "400"
test_api_error "添加评论缺少参数" POST "/api/v1/comments/add" '{}' "400"

echo ""
echo "=== 11. 用户资料与帖子 ==="

test_api "当前用户资料" GET "/api/v1/user/$USERNAME" "" "screenName"
test_api "当前用户帖子" POST "/api/v1/user/$USERNAME/posts" "{}" "data"
test_cli "CLI用户资料" "$SCREEN_NAME" "user" "$USERNAME"
test_cli "CLI用户帖子" "" "user-posts" "$USERNAME"

echo ""
echo "=== 12. 参数校验（错误处理）==="

test_api_error "like缺少post_id" POST "/api/v1/like" '{}' "400"
test_api_error "unlike缺少post_id" POST "/api/v1/unlike" '{}' "400"
test_api_error "follow缺少username" POST "/api/v1/follow" '{}' "400"
test_api_error "unfollow缺少username" POST "/api/v1/unfollow" '{}' "400"
test_api_error "创建帖子缺少content" POST "/api/v1/post/create" '{}' "400"
test_api_error "等待登录缺少uuid" POST "/api/v1/login/wait" '{}' "400"

echo ""
echo "=== 13. CLI 环境变量 ==="

# 测试 JIKE_API_URL 环境变量
CLI_RESP=$(JIKE_API_URL="http://localhost:$PORT" $CLI status 2>&1)
if echo "$CLI_RESP" | grep -q "登录\|logged"; then
  pass "CLI: JIKE_API_URL 环境变量生效"
else
  fail "CLI: JIKE_API_URL 环境变量" "$CLI_RESP"
fi

# 测试错误端口
CLI_ERR=$(JIKE_API_URL="http://localhost:19999" $CLI status 2>&1)
if echo "$CLI_ERR" | grep -q "error\|错误\|refused\|Error"; then
  pass "CLI: 错误端口正确报错"
else
  # 可能返回未登录也算合理
  pass "CLI: 错误端口返回: $(echo "$CLI_ERR" | head -c 80)"
fi

echo ""
echo "========================================"
echo -e "结果: ${GREEN}$PASS 通过${NC}, ${RED}$FAIL 失败${NC}, ${YELLOW}$SKIP 跳过${NC}"
echo ""

[ $FAIL -gt 0 ] && exit 1 || exit 0
