---
name: jike-login
description: |
  即刻扫码登录。当用户需要登录即刻、连接即刻账号、扫码时使用。
---

# 即刻登录

`P` 代表 `python3 <SKILL_DIR>/scripts/jike_client.py`，`SKILL_DIR` 取 `~/.claude/skills/jike` 或 `~/.openclaw/skills/jike` 中实际存在的路径。

## 登录流程

1. 运行 `P status` 检查是否已登录
2. 如已登录，告知用户当前账号信息
3. 如未登录：
   a. 运行 `P qrcode`，记下输出中的 `UUID` 和二维码保存路径
   b. **发送二维码图片给用户：**
      ```shell
      openclaw message send --channel telegram --target <sender_id> --media <二维码路径> --message "请用即刻 App 扫码登录"
      ```
      其中 `<sender_id>` 从会话的 Conversation info 中获取。
   c. **【关键】发送二维码后，必须立即运行 `P wait <uuid>`。不要等用户回复，发送二维码的同一轮就要执行。** 此命令会轮询最多 180 秒，扫码确认后自动返回。
   d. 成功后告知用户已登录

⚠️ **常见错误：发送二维码后停下来等用户回复，导致 `wait` 从未执行，登录永远无法完成。务必在发送二维码后立即执行 wait。**

## 登出

运行 `P status` 确认已登录后，删除 `<SKILL_DIR>/tokens.json`。

## 技术说明

- 二维码有效期约 3 分钟，超时需重新获取
- Token 保存在本地 `tokens.json`，下次启动无需重新登录
- 如遇 401 错误，token 会自动刷新；刷新失败需重新扫码
