---
name: jike-interact
description: |
  即刻社交互动：点赞、取消点赞、关注用户、取消关注。当用户想点赞、取赞、关注、取关时使用。
---

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 ~/.openclaw/skills/jike/scripts/jike_client.py`。

# 命令

| 功能 | 命令 |
|------|------|
| 点赞 | `P like <post_id>` |
| 取消点赞 | `P unlike <post_id>` |
| 关注 | `P follow <username>` |
| 取关 | `P unfollow <username>` |

# 注意

- 关注和取关操作前**必须**向用户确认
- 点赞可以直接执行，无需确认
