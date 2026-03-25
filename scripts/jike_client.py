#!/usr/bin/env python3
"""Jike MCP Server CLI Client.

Usage:
    python jike_client.py status              # Check login status
    python jike_client.py qrcode              # Get login QR code (saves to /tmp/jike-qr.png)
    python jike_client.py wait <uuid>         # Wait for QR scan confirmation
    python jike_client.py following           # Get following feeds
    python jike_client.py recommend           # Get recommended feeds
    python jike_client.py search <keyword>    # Search posts/users/topics
    python jike_client.py post-detail <id>    # Get post detail
    python jike_client.py comments <id>       # Get post comments
    python jike_client.py user <username>     # Get user profile
    python jike_client.py user-posts <username>  # Get user's posts
    python jike_client.py create-post <content>  # Create a new post
    python jike_client.py comment <post_id> <content>  # Add comment
    python jike_client.py like <post_id>      # Like a post
    python jike_client.py unlike <post_id>    # Unlike a post
    python jike_client.py follow <username>   # Follow user
    python jike_client.py unfollow <username> # Unfollow user
"""

import sys
import json
import base64
import urllib.request
import urllib.error

BASE_URL = "http://localhost:8080"


def api(method, path, data=None):
    url = BASE_URL + path
    headers = {"Content-Type": "application/json", "Accept": "application/json"}
    body = json.dumps(data).encode() if data else None
    req = urllib.request.Request(url, data=body, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return json.loads(resp.read())
    except urllib.error.HTTPError as e:
        return {"error": f"HTTP {e.code}: {e.read().decode()[:200]}"}
    except Exception as e:
        return {"error": str(e)}


def pp(obj):
    print(json.dumps(obj, indent=2, ensure_ascii=False))


def cmd_status():
    r = api("GET", "/api/v1/status")
    if r.get("logged_in"):
        u = r["user"]
        print(f"已登录: {u.get('screenName', '')} (@{u.get('username', '')})")
        s = u.get("statsCount", {})
        if isinstance(s, dict):
            print(f"  关注: {s.get('followingCount', 0)}  粉丝: {s.get('followedCount', 0)}  获赞: {s.get('liked', 0)}")
    else:
        print("未登录。请运行: python jike_client.py qrcode")


def cmd_qrcode():
    # Directly call Jike API (bypass Go server entirely for login)
    import urllib.parse
    import urllib.request
    try:
        import qrcode as qrlib
    except ImportError:
        import subprocess
        subprocess.check_call([sys.executable, "-m", "pip", "install", "--break-system-packages", "-q", "qrcode", "pillow"])
        import qrcode as qrlib

    # Create session directly via Jike API
    jike_api = "https://api.ruguoapp.com"
    req = urllib.request.Request(
        f"{jike_api}/sessions.create",
        headers={"Origin": "https://web.okjike.com", "Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(req, timeout=10) as resp:
        uuid = json.loads(resp.read()).get("uuid", "")

    if not uuid:
        print("错误: 无法创建登录会话")
        return

    # Generate QR code using Python library
    scan_url = f"https://www.okjike.com/account/scan?uuid={uuid}"
    deep_link = f"jike://page.jk/web?url={urllib.parse.quote(scan_url, safe='')}&displayHeader=false&displayFooter=false"

    qr = qrlib.QRCode(error_correction=qrlib.constants.ERROR_CORRECT_M, box_size=10, border=4)
    qr.add_data(deep_link)
    qr.make(fit=True)
    img = qr.make_image(fill_color="black", back_color="white")

    qr_path = "/tmp/jike-qr.png"
    img.save(qr_path)

    import io
    buf = io.BytesIO()
    img.save(buf, format="PNG")
    qr_base64 = base64.b64encode(buf.getvalue()).decode()

    print(f"UUID: {uuid}")
    print(f"二维码已保存: {qr_path}")
    print(f"二维码 Base64 长度: {len(qr_base64)}")
    print("请用即刻 App 的「扫一扫」功能扫描二维码")
    print(f"扫码后运行: python jike_client.py wait {uuid}")


def cmd_wait(uuid):
    import time
    import os
    print("等待扫码确认（最长 180 秒）...")
    jike_api = "https://api.ruguoapp.com"
    token_path = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "tokens.json")

    deadline = time.time() + 180
    while time.time() < deadline:
        try:
            req = urllib.request.Request(
                f"{jike_api}/sessions.wait_for_confirmation?uuid={uuid}",
                headers={"Origin": "https://web.okjike.com", "Accept": "application/json"},
            )
            with urllib.request.urlopen(req, timeout=10) as resp:
                data = json.loads(resp.read())
                access_token = data.get("x-jike-access-token", "")
                refresh_token = data.get("x-jike-refresh-token", "")
                if data.get("confirmed") and access_token:
                    # Save tokens to file (for both Go server and future CLI use)
                    tokens = {"access_token": access_token, "refresh_token": refresh_token}
                    with open(token_path, "w") as f:
                        json.dump(tokens, f, indent=2)
                    print(f"登录成功! Token 已保存到 {token_path}")
                    return
        except urllib.error.HTTPError as e:
            if e.code == 400:
                pass  # Still waiting
            elif e.code == 404:
                print("会话已过期，请重新获取二维码")
                return
        except Exception:
            pass
        time.sleep(2)

    print("登录超时（180 秒），请重试")


def cmd_following():
    r = api("POST", "/api/v1/feeds/following", {})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    for p in r.get("data", []):
        u = p.get("user", {})
        print(f"[{u.get('screenName', '?')}] {p.get('content', '')[:100]}")
        print(f"  👍{p.get('likeCount', 0)} 💬{p.get('commentCount', 0)} 🔁{p.get('repostCount', 0)}")
        print()


def cmd_recommend():
    r = api("POST", "/api/v1/feeds/recommend", {})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    for p in r.get("data", []):
        u = p.get("user", {})
        print(f"[{u.get('screenName', '?')}] {p.get('content', '')[:100]}")
        print(f"  👍{p.get('likeCount', 0)} 💬{p.get('commentCount', 0)}")
        print()


def cmd_search(keyword):
    r = api("POST", "/api/v1/search", {"keyword": keyword})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    results = r.get("data", [])
    print(f"找到 {len(results)} 条结果:")
    for p in results:
        u = p.get("user", {})
        content = p.get("content", "")
        t = p.get("type", "")
        if content:
            print(f"  [{u.get('screenName', '?')}] {content[:80]}")
            print(f"    ID: {p.get('id', '')}  👍{p.get('likeCount', 0)} 💬{p.get('commentCount', 0)}")
        elif t:
            print(f"  [{t}] {u.get('screenName', p.get('id', ''))}")
        print()


def cmd_post_detail(post_id):
    r = api("POST", "/api/v1/post/detail", {"post_id": post_id})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    pp(r)


def cmd_comments(post_id):
    r = api("POST", "/api/v1/comments/list", {"target_id": post_id})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    for c in r.get("data", []):
        u = c.get("user", {})
        print(f"[{u.get('screenName', '?')}] {c.get('content', '')}")
        print(f"  👍{c.get('likeCount', 0)}  {c.get('createdAt', '')}")
        print()


def resolve_username(input_str):
    """Resolve a Jike short URL or profile URL to a username.

    Supports: https://okjk.co/xxx, https://web.okjike.com/u/xxx, plain username.
    """
    if "://" not in input_str and "okjk.co" not in input_str:
        return input_str
    # Follow redirects manually to find the username
    import urllib.parse
    current_url = input_str
    for _ in range(10):
        req = urllib.request.Request(current_url, method="HEAD")
        try:
            opener = urllib.request.build_opener(urllib.request.HTTPRedirectHandler)
            # Use a custom opener that doesn't follow redirects
            class NoRedirect(urllib.request.HTTPRedirectHandler):
                def redirect_request(self, req, fp, code, msg, headers, newurl):
                    raise urllib.error.HTTPError(newurl, code, msg, headers, fp)
            opener = urllib.request.build_opener(NoRedirect)
            opener.open(req, timeout=10)
            break  # No redirect, we're at the final URL
        except urllib.error.HTTPError as e:
            if 300 <= e.code < 400:
                current_url = e.filename if hasattr(e, 'filename') else str(e)
                # Try to extract username from URL
                for pattern in ["/u/", "/users/"]:
                    if pattern in current_url:
                        rest = current_url.split(pattern, 1)[1]
                        username = rest.split("/")[0].split("?")[0].split("#")[0]
                        if username:
                            return username
                continue
            break
        except Exception:
            break
    # Check final URL
    for pattern in ["/u/", "/users/"]:
        if pattern in current_url:
            rest = current_url.split(pattern, 1)[1]
            username = rest.split("/")[0].split("?")[0].split("#")[0]
            if username:
                return username
    print(f"警告: 无法从 URL 解析用户名，尝试原始输入: {input_str}")
    return input_str


def cmd_user(username):
    username = resolve_username(username)
    r = api("GET", f"/api/v1/user/{username}")
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print(f"{r.get('screenName', '')} (@{r.get('username', '')})")
    print(f"  简介: {r.get('briefIntro', '') or r.get('bio', '')}")
    s = r.get("statsCount", {})
    if isinstance(s, dict):
        print(f"  关注: {s.get('followingCount', 0)}  粉丝: {s.get('followedCount', 0)}  获赞: {s.get('liked', 0)}")


def cmd_user_posts(username):
    username = resolve_username(username)
    r = api("POST", f"/api/v1/user/{username}/posts", {})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    for p in r.get("data", []):
        print(f"[{p.get('createdAt', '')}] {p.get('content', '')[:100]}")
        print(f"  ID: {p.get('id', '')}  👍{p.get('likeCount', 0)} 💬{p.get('commentCount', 0)}")
        print()


def cmd_create_post(content):
    r = api("POST", "/api/v1/post/create", {"content": content})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print(f"发布成功! ID: {r.get('id', '')}")


def cmd_comment(post_id, content):
    r = api("POST", "/api/v1/comments/add", {"target_id": post_id, "content": content})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print("评论成功!")
    pp(r)


def cmd_like(post_id):
    r = api("POST", "/api/v1/like", {"post_id": post_id})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print("点赞成功!")


def cmd_unlike(post_id):
    r = api("POST", "/api/v1/unlike", {"post_id": post_id})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print("取消点赞!")


def cmd_follow(username):
    r = api("POST", "/api/v1/follow", {"username": username})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print(f"已关注 {username}")


def cmd_unfollow(username):
    r = api("POST", "/api/v1/unfollow", {"username": username})
    if "error" in r:
        print(f"错误: {r['error']}")
        return
    print(f"已取关 {username}")


def main():
    if len(sys.argv) < 2:
        print(__doc__)
        sys.exit(1)

    cmd = sys.argv[1]
    args = sys.argv[2:]

    commands = {
        "status": (cmd_status, 0),
        "qrcode": (cmd_qrcode, 0),
        "wait": (cmd_wait, 1),
        "following": (cmd_following, 0),
        "recommend": (cmd_recommend, 0),
        "search": (cmd_search, 1),
        "post-detail": (cmd_post_detail, 1),
        "comments": (cmd_comments, 1),
        "user": (cmd_user, 1),
        "user-posts": (cmd_user_posts, 1),
        "create-post": (cmd_create_post, 1),
        "comment": (cmd_comment, 2),
        "like": (cmd_like, 1),
        "unlike": (cmd_unlike, 1),
        "follow": (cmd_follow, 1),
        "unfollow": (cmd_unfollow, 1),
    }

    if cmd not in commands:
        print(f"未知命令: {cmd}")
        print(__doc__)
        sys.exit(1)

    func, nargs = commands[cmd]
    if len(args) < nargs:
        print(f"命令 '{cmd}' 需要 {nargs} 个参数")
        sys.exit(1)

    func(*args[:nargs])


if __name__ == "__main__":
    main()
