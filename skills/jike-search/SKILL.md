---
name: jike-search
description: |
  搜索即刻内容。当用户想搜索帖子、用户、话题时使用。
---

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 ~/.openclaw/skills/jike/scripts/jike_client.py`。

# 命令

```shell
P search "关键词"
```

# 展示格式

搜索结果按类型展示：
- **帖子**：作者、内容摘要、点赞数、评论数
- **用户**：昵称、用户名、简介
- **话题**：话题名称

默认展示前 10 条结果。
