---
name: jike-comments
description: |
  查看即刻帖子详情和评论，发表评论。当用户想看某条帖子、看评论、回复评论时使用。
---

# 帖子详情与评论

## 可用工具

- `get_post_detail` — 获取帖子详情
- `get_comments` — 获取帖子评论列表
- `add_comment` — 发表评论

## 查看帖子

调用 `get_post_detail`，传入 `post_id`。如果是转发帖，`post_type` 设为 `REPOST`。

## 查看评论

调用 `get_comments`，传入：
- `target_id`：帖子 ID
- `target_type`：`ORIGINAL_POST` 或 `REPOST`（默认 `ORIGINAL_POST`）
- `load_more_key`：分页用

## 发表评论

1. 确认用户要评论的帖子和评论内容
2. 发表前**必须**向用户确认
3. 调用 `add_comment`，传入 `target_id`、`target_type`、`content`

## 展示格式

评论展示：评论者昵称、内容、点赞数、时间。嵌套回复缩进展示。
