# LinkedIn and X Launch Copy

Use [`../assets/juggleim-cover-social-1200x630.jpg`](../assets/juggleim-cover-social-1200x630.jpg) with both posts.

## LinkedIn

```text
Building production chat involves much more than opening a WebSocket.

Reconnects, acknowledgements, duplicate messages, offline users, conversation state, unread counts, groups, history, push notifications, multiple devices, moderation, and storage all become infrastructure work before chat becomes useful to the product.

JuggleIM is an open-source, self-hosted messaging backend built in Go. It provides private chat, groups, chatrooms, message history, multi-device synchronization, push integrations, REST APIs, and SDKs for major client platforms.

The community edition starts as one modular Go process plus MySQL:

git clone https://github.com/juggleim/im-server.git
cd im-server
docker compose up -d

Read the architecture and try it locally:
https://github.com/juggleim/im-server?utm_source=linkedin&utm_medium=social&utm_campaign=juggleim_oss_launch

#golang #opensource #selfhosted #websocket #realtimemessaging
```

## X — single post

```text
JuggleIM is an open-source, self-hosted messaging backend built in Go: private chat, groups, history, multi-device sync, push, REST APIs, and multi-platform SDKs.

Start it with Docker Compose ↓
https://github.com/juggleim/im-server?utm_source=x&utm_medium=social&utm_campaign=juggleim_oss_launch
```

## X — thread

```text
1/ Building production chat involves much more than opening a WebSocket: reconnects, ACKs, duplicates, offline users, history, unread state, groups, push, and multi-device sync.

2/ JuggleIM packages those concerns into an open-source, self-hosted messaging backend written in Go. Clients use Protobuf over WebSocket; business backends use REST APIs.

3/ The community deployment is one modular Go process plus MySQL, with optional MongoDB for selected message workloads. Its internal actor/RPC runtime routes work by method and target ID.

4/ Run it locally:
git clone https://github.com/juggleim/im-server.git
cd im-server
docker compose up -d

5/ Code, architecture, and docs:
https://github.com/juggleim/im-server?utm_source=x&utm_medium=social&utm_campaign=juggleim_oss_launch
```
