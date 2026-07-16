# Reddit `r/selfhosted` Draft

Check the community's current self-promotion and flair rules before posting.

## Title

```text
We built an open-source, self-hosted messaging backend in Go
```

## Body

````markdown
Hi r/selfhosted,

We have been working on [JuggleIM](https://github.com/juggleim/im-server?utm_source=reddit_selfhosted&utm_medium=community&utm_campaign=juggleim_oss_launch), an Apache-2.0 messaging backend that you can run on your own infrastructure.

It provides the infrastructure behind private chat, groups, chatrooms, history, unread state, multi-device sync, offline push, and bot conversations. Client SDKs connect using Protobuf over WebSocket, while a business backend uses REST APIs.

The community deployment is intentionally small:

- One Go server process
- MySQL 8
- Optional MongoDB for selected message workloads
- Optional object storage and push providers

You can try it with:

```bash
git clone https://github.com/juggleim/im-server.git
cd im-server
docker compose up -d
```

The local admin console is then available at `http://127.0.0.1:8090`. The documented default credentials are development-only and must be changed before exposing a deployment.

Architecture: https://github.com/juggleim/im-server/blob/master/docs/architecture.md?utm_source=reddit_selfhosted&utm_medium=community&utm_campaign=juggleim_oss_launch

For transparency, the open-source repository is a single-node server. Multi-node discovery, routing, and failover are not included in the community implementation.

We would appreciate feedback on the Docker setup, operational defaults, backup/restore expectations, and what a self-hosting troubleshooting guide should cover.
````

## Recommended flair

Use the community's project-release or self-promotion flair if one is available. Do not label it as a tutorial unless the post includes a complete tutorial.
