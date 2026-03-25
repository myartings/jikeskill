---
name: jike-post
description: |
  在即刻发布新帖子。当用户想发即刻、写动态、发帖时使用。
---

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 ~/.openclaw/skills/jike/scripts/jike_client.py`。

# 命令

```shell
P create-post "帖子内容"
```

# 流程

1. 确认用户要发布的内容
2. 发布前**必须**向用户确认
3. 运行命令发布
4. 展示发布成功的帖子信息
