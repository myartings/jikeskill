# Jike MCP Server

即刻（Jike）MCP 服务器 + OpenClaw Skills。

## 构建与运行

```bash
export PATH=$HOME/go-sdk/go/bin:$HOME/go/bin:$PATH
go build -o jike-mcp .
./jike-mcp --port 8080 --tokens tokens.json
```

Go 版本要求 >= 1.24（MCP SDK 会自动下载 Go 1.25 toolchain）。

## 架构

- 纯 REST API 调用（**不是**浏览器自动化），即刻 API 完善无需 Playwright/go-rod
- API Base URL: `https://api.ruguoapp.com`
- 认证：`x-jike-access-token` header，401 时自动用 refresh token 刷新
- Token 持久化到本地 `tokens.json`
- 双接口：MCP (`/mcp`) + REST API (`/api/v1/*`)

## 即刻 API 注意事项

这些是开发中踩过的坑，修改代码时务必注意：

### 登录流程
- QR 码内容必须是 `jike://page.jk/web?url=<encoded>&displayHeader=false&displayFooter=false`
- 内部 URL 必须是 `https://www.okjike.com/account/scan?uuid={uuid}`（**不是** `web.okjike.com/scan-login`）
- `sessions.wait_for_confirmation` 返回 400 表示等待中，200 表示已确认
- **Token 在响应体 JSON 中**（`x-jike-access-token`、`x-jike-refresh-token` 字段），不在 HTTP header 中

### 数据类型
- `avatarImage` 是嵌套对象（含 thumbnailUrl、smallPicUrl 等），不是字符串
- `statsCount` 是嵌套对象
- `createdAt` 在部分接口（如搜索）中可能为空字符串，不能用 `time.Time`，统一用 `string`
- 搜索结果 (`/1.0/search/integrate`) 包含多种类型（帖子、话题、用户），`data` 数组元素结构不统一
- 不确定类型的字段统一用 `any`，避免 JSON 反序列化失败

### API 端点
- 大部分端点是 POST，少数是 GET（如 `/1.0/users/profile?username=`）
- 公共 headers: `Origin: https://web.okjike.com`, `Content-Type: application/json`

## 项目结构

```
├── main.go              # 入口，CLI 参数
├── app_server.go        # AppServer 管理 service + MCP server
├── mcp_server.go        # 16 个 MCP tool 注册
├── mcp_handlers.go      # MCP handler 实现（解析参数、调 service、返回结果）
├── routes.go            # REST API 路由和 handlers
├── service.go           # 业务逻辑层（薄代理，调 jike client）
├── jike/                # 即刻 API 客户端
│   ├── client.go        # HTTP 客户端，自动带 token，401 自动 refresh
│   ├── login.go         # QR 码登录（create session → generate QR → poll confirm）
│   ├── types.go         # 数据类型（User, Post, Comment 等）
│   ├── feeds.go         # 动态流（关注/推荐）
│   ├── posts.go         # 帖子 CRUD
│   ├── comments.go      # 评论
│   ├── search.go        # 搜索
│   ├── user.go          # 用户资料
│   └── interactions.go  # 点赞/关注
├── tokens/tokens.go     # Token 文件读写
├── SKILL.md             # OpenClaw 根 skill
└── skills/              # 7 个 OpenClaw 子 skill
```

## 添加新 API 的流程

1. 在 `jike/` 中添加客户端方法（调 `c.Do(method, path, body)`）
2. 在 `service.go` 中添加代理方法
3. 在 `mcp_server.go` 中注册 tool（name + description + inputSchema）
4. 在 `mcp_handlers.go` 中实现 handler（parseArgs → 调 service → 返回结果）
5. 可选：在 `routes.go` 中添加 REST API endpoint

## 测试

没有自动化测试。手动测试流程：

```bash
# 启动
./jike-mcp --port 9999

# 检查状态
curl -s http://localhost:9999/api/v1/status

# 测试动态流（需要已登录，tokens.json 存在）
curl -s -X POST http://localhost:9999/api/v1/feeds/following -H "Content-Type: application/json" -d '{}'

# 测试搜索
curl -s -X POST http://localhost:9999/api/v1/search -H "Content-Type: application/json" -d '{"keyword":"AI"}'
```

## 发布

```bash
git add . && git commit -m "描述"
git push origin main
git tag v1.x.x && git push origin v1.x.x  # 触发 GitHub Actions 构建多平台二进制
```
