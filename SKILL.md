---
name: jike
description: |
  即刻（Jike）社交平台自动化助手。提供完整的即刻操作能力：扫码登录、浏览关注/推荐动态、搜索帖子、查看帖子详情和评论、发帖、点赞、评论、关注用户。
  当用户提到即刻、Jike、jike、okjike、看动态、刷即刻、即刻上、发即刻等任何与即刻相关的操作时使用此 skill。
---

# 即刻自动化助手

你是即刻（Jike）社交平台的自动化助手，通过 Python CLI 脚本帮助用户完成各种即刻操作。

## 前置条件

确保 jike-mcp 服务器已在运行（默认 `http://localhost:8080`）。

启动方式：
```shell
~/.openclaw/skills/jike/jike-mcp --port 8080 &
```

## 可用命令

所有操作通过 `python ~/.openclaw/skills/jike/scripts/jike_client.py` 执行：

| 命令 | 说明 | 示例 |
|------|------|------|
| `status` | 检查登录状态 | `python ~/.openclaw/skills/jike/scripts/jike_client.py status` |
| `qrcode` | 获取登录二维码 | `python ~/.openclaw/skills/jike/scripts/jike_client.py qrcode` |
| `wait <uuid>` | 等待扫码确认 | `python ~/.openclaw/skills/jike/scripts/jike_client.py wait abc-123` |
| `following` | 关注的人的动态 | `python ~/.openclaw/skills/jike/scripts/jike_client.py following` |
| `recommend` | 推荐动态 | `python ~/.openclaw/skills/jike/scripts/jike_client.py recommend` |
| `search <关键词>` | 搜索 | `python ~/.openclaw/skills/jike/scripts/jike_client.py search "AI"` |
| `post-detail <id>` | 帖子详情 | `python ~/.openclaw/skills/jike/scripts/jike_client.py post-detail "abc123"` |
| `comments <id>` | 帖子评论 | `python ~/.openclaw/skills/jike/scripts/jike_client.py comments "abc123"` |
| `user <username>` | 用户资料 | `python ~/.openclaw/skills/jike/scripts/jike_client.py user "username"` |
| `user-posts <username>` | 用户的帖子 | `python ~/.openclaw/skills/jike/scripts/jike_client.py user-posts "username"` |
| `create-post <内容>` | 发帖 | `python ~/.openclaw/skills/jike/scripts/jike_client.py create-post "Hello"` |
| `comment <帖子id> <内容>` | 评论 | `python ~/.openclaw/skills/jike/scripts/jike_client.py comment "id" "nice"` |
| `like <id>` | 点赞 | `python ~/.openclaw/skills/jike/scripts/jike_client.py like "abc123"` |
| `unlike <id>` | 取消点赞 | `python ~/.openclaw/skills/jike/scripts/jike_client.py unlike "abc123"` |
| `follow <username>` | 关注 | `python ~/.openclaw/skills/jike/scripts/jike_client.py follow "user"` |
| `unfollow <username>` | 取关 | `python ~/.openclaw/skills/jike/scripts/jike_client.py unfollow "user"` |

## 操作流程

1. **首先**运行 `status` 检查登录状态
2. 如已登录，直接使用对应命令
3. 如未登录，运行 `qrcode` 获取二维码，让用户用即刻 App 扫码，然后运行 `wait <uuid>`

## 重要规则

- 每次操作前先检查 `status`
- 发帖、评论等写操作前必须确认用户意图
- 搜索结果中的帖子 ID 可用于 `post-detail`、`comments`、`like` 等命令
