# Reddit `r/golang` Draft

Check the community's current project and self-promotion rules before posting.

## Title

```text
JuggleIM: a modular real-time messaging server built in Go
```

## Body

````markdown
We recently documented the architecture of [JuggleIM](https://github.com/juggleim/im-server?utm_source=reddit_golang&utm_medium=community&utm_campaign=juggleim_oss_launch), an Apache-2.0 messaging backend written in Go.

The community edition runs as one process, but keeps domain boundaries around gateways, message delivery, users, groups, conversations, history, push, files, bots, and RTC signaling.

Internally, modules register string-based actor methods such as `p_msg`, `g_msg`, and `msg_dispatch`. A Protobuf RPC envelope carries tenant, requester, target, QoS, sequence, and payload metadata. The runtime supports synchronous queries, asynchronous commands, grouped routing, and broadcast.

External clients use Protobuf over WebSocket. Trusted business backends use REST APIs. MySQL is required; MongoDB is optional for selected message and history workloads.

The full architecture and message sequence diagrams are here:

https://github.com/juggleim/im-server/blob/master/docs/architecture.md?utm_source=reddit_golang&utm_medium=community&utm_campaign=juggleim_oss_launch

Quick start:

```bash
git clone https://github.com/juggleim/im-server.git
cd im-server
docker compose up -d
```

We are preparing a reproducible benchmark harness and would value Go-specific feedback on profiling methodology, connection-load generation, actor scheduling measurements, and useful P50/P95/P99 scenarios.
````

## Discussion goal

Keep the discussion technical. Ask for review of Go architecture and benchmark methodology rather than Stars.
