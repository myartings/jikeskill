---
name: jike-login
description: |
  即刻扫码登录。当用户需要登录即刻、连接即刻账号、扫码时使用。
---

# 登录流程

1. 检查登录状态：
   ```shell
   python3 ~/.openclaw/skills/jike/scripts/jike_client.py status
   ```

2. 如已登录，告知用户当前账号信息

3. 如未登录，获取二维码：
   ```shell
   curl -s -X POST http://localhost:8080/api/v1/login/qrcode -H "Content-Type: application/json" -d '{}'
   ```
   返回 JSON 包含 `uuid` 和 `qrcode_base64`（PNG 图片的 base64 编码）。将 base64 解码为图片发送给用户，提示用即刻 App 扫码。

4. **【关键】发送二维码后，立即等待扫码确认，不要等用户回复：**
   ```shell
   python3 ~/.openclaw/skills/jike/scripts/jike_client.py wait <uuid>
   ```
   此命令会轮询最多 180 秒，扫码确认后自动返回。

5. 成功后告知用户已登录

⚠️ **常见错误：发送二维码后停下来等用户回复，导致 `wait` 从未执行，登录永远无法完成。**

# 登出

```shell
python3 ~/.openclaw/skills/jike/scripts/jike_client.py logout
```

# 技术说明

- 登录不需要公网 IP。`wait` 是客户端主动轮询即刻服务器，纯出站请求。
- 二维码有效期约 3 分钟，超时需重新获取
- 首次使用需要先运行 `cd ~/.openclaw/skills/jike && bash scripts/setup.sh` 启动服务器
