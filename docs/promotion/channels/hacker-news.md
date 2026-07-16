# Hacker News Launch Copy

Check the current Show HN rules immediately before submitting.

## Submission

**Title**

```text
Show HN: JuggleIM – an open-source messaging backend built in Go
```

**URL**

```text
https://github.com/juggleim/im-server?utm_source=hackernews&utm_medium=community&utm_campaign=juggleim_oss_launch
```

## First comment

```text
Hi HN — we built JuggleIM because production messaging quickly becomes much more than maintaining a WebSocket connection.

The open-source server is written in Go and provides private chat, groups, chatrooms, message history, unread state, multi-device synchronization, offline push integration, REST APIs, and SDKs for major client platforms. Client traffic uses Protobuf over WebSocket, while business backends use the HTTP API.

The community edition is deliberately straightforward to operate: one modular Go process plus MySQL, with optional MongoDB for selected message workloads. Internally, requests are routed through an actor/RPC runtime by method and target ID.

You can run it locally with:

git clone https://github.com/juggleim/im-server.git
cd im-server
docker compose up -d

Architecture: https://github.com/juggleim/im-server/blob/master/docs/architecture.md?utm_source=hackernews&utm_medium=community&utm_campaign=juggleim_oss_launch
Documentation: https://www.juggle.im/docs/guide/intro/?utm_source=hackernews&utm_medium=community&utm_campaign=juggleim_oss_launch

One boundary we want to state clearly: the open-source repository is a single-node implementation. Multi-node routing and failover are part of the professional offering.

We would especially value feedback on the architecture, local setup, API ergonomics, and which reproducible benchmarks would be most useful.
```

## Comment handling

- Have a maintainer available for at least two hours after submission.
- Answer the most technical questions first.
- Do not ask for upvotes.
- Do not use marketing-scale claims without benchmark links.
- Link confirmed problems to a focused GitHub Issue.
