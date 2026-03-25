---
name: jike-profile
description: |
  查看即刻用户资料和发布的帖子。当用户想看某人的主页、资料、发过的内容时使用。
---

# 查看用户资料

## 可用工具

- `get_user_profile` — 获取用户资料
- `get_user_posts` — 获取用户发布的帖子

## 查看资料

调用 `get_user_profile`，传入 `username`。

展示：昵称、用户名、简介、关注数、粉丝数、获赞数。

## 查看用户的帖子

调用 `get_user_posts`，传入 `username`，支持 `load_more_key` 分页。

## 注意事项

- 支持直接传入即刻链接（如 `https://okjk.co/xxx`），会自动解析出 username
- 如果用户提供的是昵称而非用户名，先尝试搜索找到对应的 username
- 展示资料时格式清晰，重要数据突出
