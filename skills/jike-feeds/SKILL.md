---
name: jike-feeds
description: |
  浏览即刻动态流。当用户想看关注的人的动态、推荐内容、刷即刻时使用。
---

# 浏览即刻动态

## 可用工具

- `get_following_feeds` — 获取关注的人的动态
- `get_recommend_feeds` — 获取推荐动态

## 使用方式

1. 默认展示关注动态（`get_following_feeds`）
2. 如用户要求看推荐内容，使用 `get_recommend_feeds`
3. 首次请求不传 `load_more_key`，加载更多时传入上次返回的 `loadMoreKey`

## 展示格式

每条动态展示：
- 作者昵称（screenName）和用户名（username）
- 发布时间
- 正文内容（content）
- 图片数量（如有）
- 点赞数、评论数
- 所属话题（如有）

简洁展示，每条之间用分隔线隔开。默认展示 5-10 条，用户可要求查看更多。
