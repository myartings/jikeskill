---
name: jike-feeds
description: |
  浏览即刻动态流。当用户想看关注的人的动态、推荐内容、刷即刻时使用。
---

# 前置检查

如果 `~/.openclaw/skills/jike/scripts/jike_client.py` 不存在，说明安装不完整，先执行：
```shell
rm -rf ~/.openclaw/skills/jike && git clone https://github.com/myartings/jikeskill.git ~/.openclaw/skills/jike && cd ~/.openclaw/skills/jike && bash scripts/setup.sh
```

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 ~/.openclaw/skills/jike/scripts/jike_client.py`。

# 命令

| 功能 | 命令 |
|------|------|
| 关注动态 | `P following` |
| 推荐动态 | `P recommend` |

默认展示关注动态。用户要求看推荐时用 `recommend`。

# 展示格式

每条动态展示：作者昵称、发布时间、正文、点赞数、评论数。简洁展示，每条之间用分隔线隔开。默认展示 5-10 条。
