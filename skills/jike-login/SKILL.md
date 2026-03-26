---
name: jike-login
description: |
  即刻扫码登录。当用户需要登录即刻、连接即刻账号、扫码时使用。
---

# 规则

**只用下面的 python3 命令，禁止使用 curl 或其他方式。**

`P` 代表 `python3 ~/.openclaw/skills/jike/scripts/jike_client.py`。

# 登录流程

1. 运行 `P status` 检查是否已登录
2. 如已登录，告知用户当前登录的账号信息
3. 如未登录：
   a. 运行 `P qrcode`，记下输出的 `UUID`，二维码保存在 `/tmp/jike-qr.png`
   b. 将二维码图片展示给用户，提示打开即刻 App 扫描
   c. **【关键】展示二维码后，立即运行 `P wait <uuid>`。不要等用户回复。** 此命令会轮询最多 180 秒，扫码确认后自动返回。
   d. 成功后告知用户已登录

⚠️ **常见错误：展示二维码后停下来等用户回复，导致 `wait` 从未执行，登录永远无法完成。**

# 登出

运行 `P logout`。
