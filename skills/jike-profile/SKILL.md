---
name: jike-profile
description: |
  查看即刻用户资料和发布的帖子。当用户想看某人的主页、资料、发过的内容、某个用户的动态时使用。
  支持即刻链接（okjk.co/xxx）和用户名。
---

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 <SKILL_DIR>/scripts/jike_client.py`（SKILL_DIR：`${CLAUDE_PLUGIN_ROOT}` 或 `~/.openclaw/skills/jike` 或 `~/.claude/skills/jike`，取存在的路径）。

# 命令

| 功能 | 命令 |
|------|------|
| 查看用户资料 | `P user <username 或 okjk.co链接 或短码>` |
| 查看用户帖子 | `P user-posts <username 或 okjk.co链接 或短码>` |

支持三种输入：username、即刻链接（`https://okjk.co/xxx`）、裸短码（如 `rAgUmv`），自动解析。

# 示例

```shell
P user rAgUmv
P user-posts https://okjk.co/rAgUmv
P user-posts 27BF807A-FA4D-4B01-AAFD-05FAAA674335
```

# 展示格式

资料展示：昵称、用户名、简介、关注数、粉丝数、获赞数。

如果用户提供的是昵称而非用户名，先用 `P search "昵称"` 搜索找到 username。
