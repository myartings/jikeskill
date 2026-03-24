# jike-mcp

即刻（Jike）MCP 服务器 + OpenClaw Skills，让 AI 助手帮你刷即刻。

## 功能

- 扫码登录即刻
- 浏览关注/推荐动态流
- 搜索帖子、用户、话题
- 查看帖子详情和评论
- 发帖、评论、点赞、关注

## 快速开始

### 安装为 OpenClaw Skill

```bash
# 克隆到 OpenClaw skills 目录
cd ~/.openclaw/skills/
git clone https://github.com/myartings/jikeskill.git jike

# 编译 MCP 服务器
cd jike
go build -o jike-mcp .

# 启动服务器
./jike-mcp
```

或者使用已编译的二进制（见 [Releases](https://github.com/myartings/jikeskill/releases)）：

```bash
cd ~/.openclaw/skills/
git clone https://github.com/myartings/jikeskill.git jike
cd jike

# 下载对应平台的二进制
# Linux amd64:
curl -L https://github.com/myartings/jikeskill/releases/latest/download/jike-mcp-linux-amd64 -o jike-mcp
chmod +x jike-mcp
./jike-mcp
```

### 使用 mcporter

```bash
# 先启动 MCP 服务器（同上）
# 然后添加到 OpenClaw
mcporter config add jike-mcp http://localhost:8080/mcp
```

### 手动配置

编辑 `~/.openclaw/openclaw.json`：

```json
{
  "mcpServers": {
    "jike-mcp": {
      "type": "streamableHttp",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

## 使用方式

启动 MCP 服务器后，在 OpenClaw 中对话即可：

- "帮我登录即刻" → 扫码登录
- "看看我关注的人发了什么" → 浏览关注动态
- "搜索关于 AI 的帖子" → 搜索
- "帮我发一条即刻" → 发帖
- "点赞这条帖子" → 点赞

## 命令行参数

```bash
./jike-mcp [flags]

  -port     服务器端口（默认 8080）
  -tokens   Token 存储路径（默认 tokens.json）
```

## API

MCP 服务器同时提供两种接口：

- **MCP**: `http://localhost:8080/mcp`（Streamable HTTP）
- **REST API**: `http://localhost:8080/api/v1/`

### MCP Tools

| Tool | 说明 |
|------|------|
| `check_login_status` | 检查登录状态 |
| `get_login_qrcode` | 获取登录二维码 |
| `wait_for_login` | 等待扫码确认 |
| `logout` | 登出 |
| `get_following_feeds` | 关注动态 |
| `get_recommend_feeds` | 推荐动态 |
| `search` | 搜索 |
| `get_post_detail` | 帖子详情 |
| `get_comments` | 评论列表 |
| `create_post` | 发帖 |
| `add_comment` | 评论 |
| `get_user_profile` | 用户资料 |
| `get_user_posts` | 用户帖子 |
| `like_post` / `unlike_post` | 点赞/取赞 |
| `follow_user` / `unfollow_user` | 关注/取关 |

## 构建

```bash
go build -o jike-mcp .
```

## License

MIT
