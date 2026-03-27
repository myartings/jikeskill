---
name: jike-search
description: |
  搜索即刻内容，浏览圈子帖子。当用户想搜索帖子、用户、话题、圈子，或查看某个圈子的内容时使用。
---

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 <SKILL_DIR>/scripts/jike_client.py`（SKILL_DIR：`${CLAUDE_PLUGIN_ROOT}` 或 `~/.openclaw/skills/jike` 或 `~/.claude/skills/jike`，取存在的路径）。

# 命令

| 功能 | 命令 |
|------|------|
| 搜索 | `P search "关键词"` |
| 查看圈子帖子 | `P topic-feed <topic_id>` |

搜索结果包含帖子、用户和圈子。圈子结果会显示 `ID`，用这个 ID 可以用 `topic-feed` 查看圈子内的帖子。

# 流程

查看某个圈子的内容：
1. 先用 `P search "圈子名"` 搜索，找到圈子的 ID
2. 再用 `P topic-feed <topic_id>` 查看圈子内的帖子
