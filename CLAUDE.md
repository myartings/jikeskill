# Jike MCP Server

即刻（Jike）MCP 服务器，提供 MCP 协议和 REST API 两种接口。

## 构建与运行

```bash
go build -o jike-mcp .
./jike-mcp --port 8080 --tokens tokens.json
```

## 架构

- 纯 REST API 调用，无浏览器自动化
- API Base URL: `https://api.ruguoapp.com`
- 认证方式：QR 码扫码登录，Token 持久化到本地 JSON 文件
- MCP endpoint: `/mcp`（Streamable HTTP）
- REST API: `/api/v1/*`

## 项目结构

- `jike/` - 即刻 API 客户端（client, login, feeds, posts, comments, search, user, interactions）
- `tokens/` - Token 持久化
- `main.go` - 入口
- `app_server.go` - AppServer
- `mcp_server.go` - MCP tool 注册
- `mcp_handlers.go` - MCP handler 实现
- `routes.go` - HTTP 路由和 REST API handlers
- `service.go` - 业务逻辑层

## 登录流程

1. 调用 `get_login_qrcode` 获取二维码
2. 用即刻 App 扫描二维码
3. 调用 `wait_for_login` 等待确认（传入 UUID）

## 依赖

- `github.com/modelcontextprotocol/go-sdk` - MCP SDK
- `github.com/gin-gonic/gin` - HTTP 框架
- `github.com/skip2/go-qrcode` - QR 码生成
