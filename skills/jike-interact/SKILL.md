---
name: jike-interact
description: |
  即刻社交互动：点赞、取消点赞、关注用户、取消关注。当用户想点赞、取赞、关注、取关时使用。
---

# 即刻社交互动

## 可用工具

- `like_post` — 点赞帖子
- `unlike_post` — 取消点赞
- `follow_user` — 关注用户
- `unfollow_user` — 取消关注

## 点赞/取消点赞

调用 `like_post` 或 `unlike_post`，传入：
- `post_id`：帖子 ID
- `target_type`：`ORIGINAL_POST` 或 `REPOST`（默认 `ORIGINAL_POST`）

## 关注/取消关注

调用 `follow_user` 或 `unfollow_user`，传入 `username`。

## 注意事项

- 关注和取关操作前必须向用户确认
- 点赞操作可以直接执行，无需确认
- 操作完成后告知用户结果
