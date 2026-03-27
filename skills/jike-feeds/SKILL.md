---
name: jike-feeds
description: |
  浏览即刻动态流和热门内容。当用户想看关注的人的动态、推荐内容、热门帖子、刷即刻时使用。
---

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 <SKILL_DIR>/scripts/jike_client.py`（SKILL_DIR：`${CLAUDE_PLUGIN_ROOT}` 或 `~/.openclaw/skills/jike` 或 `~/.claude/skills/jike`，取存在的路径）。

# 命令

| 功能 | 命令 |
|------|------|
| 关注动态 | `P following` |
| 推荐动态 | `P recommend` |
| 热门帖子 | `P hot`（推荐动态按点赞排序） |

默认展示关注动态。用户要求看推荐时用 `recommend`，看热门/热榜时用 `hot`。

# 展示格式

每条动态展示：作者昵称、发布时间、正文、点赞数、评论数。简洁展示，每条之间用分隔线隔开。默认展示 5-10 条。
