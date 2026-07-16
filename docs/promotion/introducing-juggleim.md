---
title: "JuggleIM: An Open-Source, Self-Hosted Messaging Backend Built in Go"
description: "Add private chat, groups, chatrooms, multi-device sync, push notifications, and real-time messaging to your product without building the infrastructure from scratch."
tags: go, opensource, websocket, selfhosted
cover_image: https://raw.githubusercontent.com/juggleim/im-server/master/docs/promotion/assets/juggleim-cover-devto.jpg
canonical_url: https://dev.to/yuwnloyblog/juggleim-an-open-source-self-hosted-messaging-backend-built-in-go-43nh
published: true
---

Real-time messaging looks simple until you try to build it.

A prototype can send a message over a WebSocket in an afternoon. A production messaging system must also handle reconnects, acknowledgements, duplicate messages, offline users, conversation state, unread counts, groups, message history, push notifications, multiple devices, moderation, storage, and operational visibility.

That is a large amount of infrastructure to build before chat becomes useful to your actual product.

[JuggleIM](https://github.com/juggleim/im-server?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch) is an open-source, self-hosted messaging backend designed to provide that foundation. It is written in Go, uses Protobuf over WebSocket for client connections, exposes REST APIs for business services, and includes SDKs for major client platforms.

At the time of writing, the project has earned more than 3,500 GitHub stars and 360 forks. More importantly, it can be started locally with one Docker Compose command.

## What can you build with it?

JuggleIM is not a complete social product or customer-support application. It is the messaging infrastructure underneath those products.

Typical use cases include:

- Private and group chat inside a mobile or web application
- Customer-service messaging
- Marketplace conversations between buyers and sellers
- Live-stream and community chatrooms
- Internal collaboration tools
- Device-to-cloud messaging for IoT products
- Real-time conversations with AI assistants and bots
- Multi-tenant communication features for SaaS products

The server provides the messaging layer while your business backend continues to own product-specific concepts such as accounts, permissions, subscriptions, orders, or CRM records.

## Why not build it from scratch?

Sending bytes is only the beginning. A useful messaging system needs answers to questions such as:

- What happens when a client reconnects after losing its network?
- How do you prevent a retried publication from creating duplicate messages?
- How are messages synchronized across a user's phone, browser, and desktop client?
- How do you maintain unread counts and conversation ordering?
- When should an offline push notification be sent?
- How do group membership, mute settings, blocks, and moderation affect delivery?
- Where do message history and delivery state live?
- How does one deployment isolate multiple applications or tenants?

JuggleIM already models these concerns. The server includes message IDs and sequence numbers, QoS-aware acknowledgements, duplicate-publication filtering, sendbox and history paths, conversation state, group membership, online presence, push routing, and tenant-scoped credentials.

This lets a team spend more time on its product experience and less time rebuilding messaging plumbing.

## A modular Go architecture

The open-source server runs as a single Go process with clear internal module boundaries.

![JuggleIM system architecture](https://raw.githubusercontent.com/juggleim/im-server/master/docs/promotion/assets/system-overview.png)

There are four primary entry points:

| Entry point        | Default port | Purpose                                                    |
| ------------------ | -----------: | ---------------------------------------------------------- |
| Server API Gateway |       `9001` | REST APIs used by your trusted business backend            |
| Navigator          |       `9002` | Validates client tokens and returns the WebSocket endpoint |
| Connect Manager    |       `9003` | Maintains Protobuf-over-WebSocket client connections       |
| Admin Gateway      |       `8090` | Serves the administration console and APIs                 |

Requests from these gateways enter an internal actor and RPC runtime. Domain modules then handle messaging, users, friends, presence, conversations, groups, history, push, file storage, subscriptions, bots, moderation, and RTC room signaling.

Actor methods are routed by a stable target ID such as a user, group, or conversation. The runtime supports synchronous queries, asynchronous commands, grouped routing, and broadcast. This keeps service boundaries explicit without requiring the operational overhead of a distributed microservice deployment for the community edition.

You can read the full [English architecture guide](https://github.com/juggleim/im-server/blob/master/docs/architecture.md?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch) or its [Chinese version](https://github.com/juggleim/im-server/blob/master/docs/architecture_zh.md?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch).

## Storage without unnecessary complexity

The default deployment requires MySQL 8. It stores application configuration and the core domain data needed by the messaging system.

MongoDB is optional. When configured as the message storage engine, it provides alternative collections for message, history, and push workloads. A local LevelDB-backed KV store is also available for internal timestamp-ordered data.

Attachments can integrate with S3-compatible storage, MinIO, Alibaba Cloud OSS, or Qiniu. Offline notifications can be delivered through APNs, FCM, and supported Android vendor channels.

The result is a small starting topology that can grow with the product instead of requiring a long list of infrastructure dependencies on day one.

## Start it locally

The fastest way to try JuggleIM is Docker Compose:

```bash
git clone https://github.com/juggleim/im-server.git
cd im-server
docker compose up -d
```

After the containers become healthy, open the admin console:

```text
http://127.0.0.1:8090
```

The development credentials are:

```text
username: admin
password: 123456
```

These credentials are intended only for local development and must be changed before a production deployment.

Create the first tenant from the admin API:

```bash
curl --request POST \
  --url http://127.0.0.1:8090/admingateway/apps/create \
  --header 'Content-Type: application/json' \
  --data '{"app_key":"my-app","app_name":"My App"}'
```

The response contains an `app_key` and `app_secret`. The secret belongs only on your trusted business backend; it must never be embedded in a client application.

From there, the business backend can use the server API on port `9001`, while client SDKs discover and connect to the WebSocket service on port `9003`.

The complete setup instructions are available in the [deployment guide](https://www.juggle.im/docs/guide/deploy/quickdeploy/?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch).

## More than one client platform

A messaging backend is only useful if clients can integrate with it. The JuggleIM organization provides SDKs and demo applications across the ecosystem, including:

- Android
- iOS
- Web
- React Native
- Flutter
- HarmonyOS
- Desktop applications
- Go and Java server-side SDKs

There are also demo business services and web clients that show how authentication, users, friends, groups, and application-specific workflows can sit above the core IM server.

Browse the complete ecosystem from the [JuggleIM GitHub organization](https://github.com/juggleim?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch).

## Multi-tenancy is part of the design

JuggleIM carries an `app_key` through API requests, internal RPC calls, storage operations, and message delivery. This allows one deployment to host multiple isolated applications.

That is useful for:

- SaaS products that provision a separate messaging space per customer
- Teams running staging and production applications on shared infrastructure
- Platforms operating several brands or regional applications
- Developers building a reusable communication service for multiple products

Tenant isolation is not something added as an afterthought; it is part of the request and routing context.

## A transparent community-edition boundary

The open-source repository implements a modular single-node server. Its internal package is named `gmicro.Cluster`, but the community implementation routes work to the current node.

Multi-node discovery, cross-node routing, failover, horizontal scaling, and commercial support belong to the professional offering. We prefer to make this boundary explicit rather than imply that the community repository provides behavior it does not contain.

For many product teams, a self-hosted single-node deployment is an effective way to develop, validate, and operate an initial messaging workload. Teams that later need a distributed topology can evaluate the professional edition against their measured requirements.

## Security starts with clear boundaries

JuggleIM uses tenant-scoped credentials and token-based client authentication. In a production environment, HTTP and WebSocket traffic should be protected with HTTPS and WSS through a trusted TLS termination layer.

Operators should also restrict the admin console, diagnostics endpoint, database ports, log uploads, and storage credentials. Default development credentials should never remain enabled in production.

The project does not claim automatic end-to-end encryption between chat participants. If your product requires application-layer content encryption, evaluate it separately as part of the client and key-management design.

Clear security claims are more valuable than broad promises.

## Where the project is going

The project is actively improving its developer experience and technical transparency. Current community work includes:

- Reproducible performance benchmarks
- Broader CI and automated security checks
- More complete configuration references
- End-to-end API examples
- Docker Compose troubleshooting documentation

These tasks are tracked publicly in [GitHub Issues](https://github.com/juggleim/im-server/issues?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch), including issues labeled [`good first issue`](https://github.com/juggleim/im-server/labels/good%20first%20issue?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch) for new contributors.

## Try it, inspect it, and tell us what is missing

If you are adding messaging to a product, JuggleIM gives you a practical starting point that you can run, inspect, and self-host.

- [Star or fork JuggleIM on GitHub](https://github.com/juggleim/im-server?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch)
- [Read the documentation](https://www.juggle.im/docs/guide/intro/?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch)
- [Follow the quick deployment guide](https://www.juggle.im/docs/guide/deploy/quickdeploy/?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch)
- [Inspect the architecture](https://github.com/juggleim/im-server/blob/master/docs/architecture.md?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch)
- [Ask a question or share an idea](https://github.com/juggleim/im-server/discussions?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch)

If the project solves a problem for you, a GitHub star helps other developers discover it. If it does not yet fit your use case, open a discussion and describe what you are building. Concrete feedback is how open-source infrastructure becomes more useful.
