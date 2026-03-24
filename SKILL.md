---
name: jike
description: |
  即刻（Jike）社交平台自动化助手。提供完整的即刻操作能力：扫码登录、浏览关注/推荐动态、搜索帖子、查看帖子详情和评论、发帖、点赞、评论、关注用户。
  当用户提到即刻、Jike、jike、okjike、看动态、刷即刻、即刻上、发即刻等任何与即刻相关的操作时使用此 skill。
  必须且只能通过 scripts/jike_client.py 脚本来执行所有即刻操作，不要尝试其他方式。
---

# 即刻自动化助手

通过 Python 脚本操作即刻平台。**所有操作必须且只能通过下面的命令执行，禁止使用 curl 或其他方式直接调用 API。**

## 首次使用

首次使用前需要初始化：

```shell
cd ~/.openclaw/skills/jike && bash scripts/setup.sh
```

这会自动编译服务器并启动。之后无需重复执行。

## 所有可用命令

脚本路径：`~/.openclaw/skills/jike/scripts/jike_client.py`

### 认证

```shell
# 检查登录状态
python3 ~/.openclaw/skills/jike/scripts/jike_client.py status

# 获取登录二维码（保存到 /tmp/jike-qr.png）
python3 ~/.openclaw/skills/jike/scripts/jike_client.py qrcode

# 等待扫码确认（uuid 来自 qrcode 命令的输出）
python3 ~/.openclaw/skills/jike/scripts/jike_client.py wait <uuid>
```

### 浏览

```shell
# 关注的人的动态
python3 ~/.openclaw/skills/jike/scripts/jike_client.py following

# 推荐动态
python3 ~/.openclaw/skills/jike/scripts/jike_client.py recommend
```

### 搜索

```shell
# 搜索帖子/用户/话题
python3 ~/.openclaw/skills/jike/scripts/jike_client.py search "关键词"
```

### 帖子详情与评论

```shell
# 查看帖子详情（id 来自搜索或动态流的结果）
python3 ~/.openclaw/skills/jike/scripts/jike_client.py post-detail <post_id>

# 查看帖子的评论列表
python3 ~/.openclaw/skills/jike/scripts/jike_client.py comments <post_id>
```

### 用户

```shell
# 查看用户资料
python3 ~/.openclaw/skills/jike/scripts/jike_client.py user <username>

# 查看用户发布的帖子
python3 ~/.openclaw/skills/jike/scripts/jike_client.py user-posts <username>
```

### 互动

```shell
# 发帖
python3 ~/.openclaw/skills/jike/scripts/jike_client.py create-post "帖子内容"

# 评论
python3 ~/.openclaw/skills/jike/scripts/jike_client.py comment <post_id> "评论内容"

# 点赞 / 取消点赞
python3 ~/.openclaw/skills/jike/scripts/jike_client.py like <post_id>
python3 ~/.openclaw/skills/jike/scripts/jike_client.py unlike <post_id>

# 关注 / 取关
python3 ~/.openclaw/skills/jike/scripts/jike_client.py follow <username>
python3 ~/.openclaw/skills/jike/scripts/jike_client.py unfollow <username>
```

## 操作流程

1. **每次操作前**先运行 `status` 检查登录状态
2. 如未登录：运行 `qrcode` → 将二维码图片发给用户 → 用户用即刻 App 扫码 → 运行 `wait <uuid>`
3. 如已登录：直接执行对应命令
4. 搜索/动态结果中的帖子 `ID` 可用于 `post-detail`、`comments`、`like` 等命令
5. 当用户想看某条帖子的评论时，使用 `comments <post_id>` 命令

## 重要规则

- **只使用上面列出的 python3 命令，不要用 curl 或其他方式**
- 发帖、评论、关注等写操作前必须确认用户意图
- 二维码图片路径固定为 `/tmp/jike-qr.png`
