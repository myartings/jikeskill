---
name: jike-comments
description: |
  查看即刻帖子详情和评论，发表评论。当用户想看某条帖子、看评论、回复评论时使用。
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
| 帖子详情 | `P post-detail <post_id>` |
| 评论列表 | `P comments <post_id>` |
| 发表评论 | `P comment <post_id> "评论内容"` |

post_id 来自搜索或动态流的结果。

# 注意

- 发表评论前**必须**向用户确认内容
