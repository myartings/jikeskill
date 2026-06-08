---
name: jike
description: "即刻（Jike）社交平台助手。浏览动态、搜索、查看用户、发帖、评论、点赞、关注。当用户提到即刻、Jike、okjike、刷即刻、发即刻时使用。"
version: 0.1.0
author: myartings
license: MIT
metadata:
  tags: [jike, social, cn, okjike]
---

# 即刻 (Jike)

## 规则

1. **只用下面的 python3 命令。禁止用 curl、wget、httpie 或任何其他方式。**
2. Skill 目录（`SKILL_DIR`）：优先取实际存在的路径，依次检查 `~/.claude/skills/jike`、`~/.openclaw/skills/jike`、`~/workspace/myartings-agent-skills/skills/jike`
3. 首次使用先运行初始化：`cd <SKILL_DIR> && bash scripts/setup.sh`
4. 每次操作前先运行 `status` 检查登录状态
5. 发帖、评论、关注等写操作前必须确认用户意图

## 命令

以下是全部可用命令，`P` 代表 `python3 <SKILL_DIR>/scripts/jike_client.py`。

| 功能 | 命令 |
|------|------|
| 检查登录 | `P status` |
| 登录二维码 | `P qrcode`（二维码保存到 /tmp/jike-qr.png） |
| 等待扫码 | `P wait <uuid>`（uuid 来自 qrcode 输出，展示二维码后立即执行） |
| 关注动态 | `P following` |
| 推荐动态 | `P recommend` |
| 热门帖子 | `P hot`（推荐动态按点赞排序） |
| 搜索 | `P search "关键词"`（结果含帖子、用户、圈子） |
| 圈子帖子 | `P topic-feed <topic_id>`（topic_id 来自搜索结果） |
| 帖子详情 | `P post-detail <post_id>` |
| 评论列表 | `P comments <post_id>` |
| 查看用户 | `P user <username 或 okjk.co链接 或短码>` |
| 用户帖子 | `P user-posts <username 或 okjk.co链接 或短码>` |
| 发帖 | `P create-post "内容"` |
| 评论 | `P comment <post_id> "内容"` |
| 点赞 | `P like <post_id>` |
| 取消点赞 | `P unlike <post_id>` |
| 关注 | `P follow <username>` |
| 取关 | `P unfollow <username>` |

## 示例

```shell
python3 <SKILL_DIR>/scripts/jike_client.py status
python3 <SKILL_DIR>/scripts/jike_client.py user-posts https://okjk.co/rAgUmv
python3 <SKILL_DIR>/scripts/jike_client.py user rAgUmv
python3 <SKILL_DIR>/scripts/jike_client.py search "AI"
```
