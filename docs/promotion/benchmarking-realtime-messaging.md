---
title: "Benchmarking Real-Time Messaging Without Mixing Up ACKs, Deliveries, and Fan-Out"
description: "How we built a reproducible Go benchmark for private chat and group messaging, with separate connection, ACK, and end-to-end delivery metrics."
tags: go, performance, websocket, opensource
cover_image: https://raw.githubusercontent.com/juggleim/im-server/master/docs/promotion/assets/juggleim-cover-devto.jpg
canonical_url: https://dev.to/yuwnloyblog/benchmarking-real-time-messaging-without-mixing-up-acks-deliveries-and-fan-out-1n6f
published: true
---

A messaging benchmark can report an impressive number and still tell you almost nothing.

“Messages per second” might mean WebSocket frames written by a load generator, publishes accepted by a server, messages committed to storage, or callbacks observed by recipients. In group chat, one accepted publish may produce hundreds or thousands of deliveries. Connection storms and steady-state delivery exercise different parts of a system, yet they are often rolled into one average.

We ran into this problem while improving performance transparency for [JuggleIM](https://github.com/juggleim/im-server?utm_source=devto&utm_medium=article&utm_campaign=juggleim_benchmark), an open-source messaging server written in Go. Before publishing bigger numbers, we wanted a benchmark whose semantics were explicit and whose results could be reproduced.

The result is a small Go harness for private and group chat that reports connection setup, server acknowledgements, and recipient delivery separately. This article explains the decisions behind it, the mistakes it tries to avoid, and how to run it yourself.

## Start by defining what a message means

For a real-time messaging system, at least four events matter:

1. **Connection established:** the client completed the production authentication and WebSocket handshake.
2. **Publish attempted:** the load generator handed a message to a connected client.
3. **Publish acknowledged:** the server accepted the publish and returned its protocol-level ACK.
4. **Delivery observed:** a recipient client decoded the message and invoked its application callback.

These events answer different questions.

Connection latency tells you how quickly a deployment can admit clients. ACK latency covers client serialization, transport, server processing, and the ACK return path. Delivery latency covers the full sender-to-recipient path. Delivery count reveals fan-out.

Combining them into one number hides the system behavior we actually want to understand.

## Two workloads, two different shapes

The harness implements two workloads using JuggleIM's production Protobuf-over-WebSocket protocol.

### Private chat: a ring

Every connected client is a sender. Clients form a ring:

```text
client 1 -> client 2
client 2 -> client 3
...
client N -> client 1
```

The configured rate is the aggregate publish rate across the ring. This spreads work across clients without generating a quadratic number of conversations or concentrating all sends on one connection.

For every measured publish, the harness records:

- publish-to-ACK latency;
- whether the ACK succeeded or failed;
- publish-to-recipient-callback latency;
- the number of recipient callbacks observed.

### Group chat: explicit fan-out

All clients join one synthetic group. A configurable subset acts as senders while every connected member can receive messages.

The report does not call recipient callbacks “message throughput.” It reports two rates:

- **ACK throughput:** accepted group publishes per second;
- **delivery throughput:** callbacks observed across group members per second.

If 100 group members receive one accepted publish, that is one successful ACK and roughly 100 delivery events. Both values are useful. Treating them as interchangeable is not.

## Measure connection establishment separately

Opening thousands of authenticated WebSocket connections exercises token validation, connection state, file descriptors, goroutine scheduling, and network setup. Steady-state messaging exercises routing, acknowledgements, persistence, conversation state, and recipient fan-out.

The harness therefore completes and records all connection attempts before warm-up begins. It refuses to publish a partial baseline if any configured client fails to connect.

Connection metrics include attempts, successes, failures, error codes, and P50/P95/P99 latency. They are never mixed into steady-state message latency.

## Warm up before recording

The first few operations of a process may include lazy initialization, empty caches, new database connections, or one-time allocations. Measuring them together with steady state makes short tests especially noisy.

Each workload has three phases:

```text
setup -> connect -> warm up -> measure -> delivery grace
```

Warm-up messages use the real protocol and storage behavior but are tagged as `warmup` and excluded from the report. Measured messages carry a run ID, phase, sequence number, nanosecond send timestamp, and enough padding to reach the requested payload size.

After sending stops, a configurable grace period collects deliveries already in flight.

## Use production semantics, not a convenient mock

The benchmark creates synthetic users through JuggleIM's signed server API, obtains real client tokens, and connects through the same WebSocket protocol used by client SDKs. Group workloads create actual group membership before clients connect.

Stored and counted text messages are enabled by default. Turning persistence off would produce a different and usually much easier workload, so the setting is included in every report.

The benchmark application secret is read only from an environment variable. It is not placed in process arguments or written into result files.

## Percentiles are more useful than one average

An average can look healthy while a meaningful fraction of users experience long delays. The report includes minimum, mean, P50, P95, P99, and maximum latency for:

- connection establishment;
- successful publish acknowledgements;
- observed deliveries.

Failures are counted by error code instead of being silently discarded or folded into successful latency.

The JSON report also records the workload and environment:

- server commit and dirty-tree state;
- client count and group sender count;
- target rate, warm-up, duration, and delivery grace;
- payload size and persistence flags;
- OS, architecture, Go version, CPU model, logical CPU count, and memory;
- database and deployment labels.

Without that context, two benchmark files are not meaningfully comparable.

## A verified smoke run

We checked the full workflow against the repository's local Docker Compose stack. The server and load generator shared one Apple M2 host, the measurement lasted only ten seconds, and the working tree contained the new harness.

That makes this a functional smoke baseline, **not a production capacity claim**.

The workload used 20 connected clients, stored and counted 256-byte messages, a three-second warm-up, a ten-second measurement, and a target of 50 publishes per second.

| Scenario | Connections | Successful ACKs | ACK P95 / P99 | Observed deliveries | Delivery P95 / P99 |
| --- | ---: | ---: | ---: | ---: | ---: |
| Private ring | 20/20 | 499/499 | 2.37 / 17.09 ms | 499 | 8.47 / 36.75 ms |
| 20-member group, 2 senders | 20/20 | 499/499 | 4.54 / 30.19 ms | 9,492 | 35.29 / 115.88 ms |

The group workload makes the distinction visible: 499 accepted publishes created 9,492 observed client deliveries. Saying the system handled either “49.9 messages per second” or “949.2 messages per second” without naming the metric would be misleading.

The full [private-chat JSON](https://github.com/juggleim/im-server/blob/master/docs/benchmarks/results/local-smoke-20260720-private.json?utm_source=devto&utm_medium=article&utm_campaign=juggleim_benchmark) and [group-chat JSON](https://github.com/juggleim/im-server/blob/master/docs/benchmarks/results/local-smoke-20260720-group.json?utm_source=devto&utm_medium=article&utm_campaign=juggleim_benchmark) are public.

## Run both workloads locally

The one-command runner requires Docker Compose, Go, `curl`, `jq`, `openssl`, and `shasum`:

```bash
git clone https://github.com/juggleim/im-server.git
cd im-server
scripts/run-benchmark.sh all
```

It builds and starts the isolated local stack, creates a benchmark application, registers synthetic users, executes both workloads, and writes JSON reports under `benchmark-results/`.

The defaults are intentionally modest. Override them with environment variables:

```bash
BENCH_CLIENTS=500 \
BENCH_RATE=1000 \
BENCH_WARMUP=30s \
BENCH_DURATION=2m \
BENCH_PAYLOAD_BYTES=1024 \
scripts/run-benchmark.sh all
```

The CLI rejects non-loopback targets by default. A separate `--allow-non-loopback` flag exists for dedicated benchmark environments, but the runner should never be pointed at shared staging or production infrastructure.

You can inspect the [complete benchmark methodology and CLI reference](https://github.com/juggleim/im-server/blob/master/BENCHMARKS.md?utm_source=devto&utm_medium=article&utm_campaign=juggleim_benchmark) before running it.

## What a serious capacity test still needs

A reproducible harness is necessary, but it is not sufficient for a defensible capacity claim.

For a publishable performance study, we would additionally:

1. Run the server, database, and load generator on separate named hosts.
2. Use a clean server commit and preserve every result file.
3. Record container limits, MySQL configuration, storage type, network topology, and TLS status.
4. Run multiple trials after a longer warm-up and publish all trials, not only the best one.
5. Collect server CPU, memory, network, disk, database, and Go profile data.
6. Increase load until errors or latency show the actual saturation point.
7. Test private chat and multiple group sizes independently.
8. Compare systems only when persistence, acknowledgement, payload, fan-out, and delivery semantics are equivalent.

Large numbers without these controls are demonstrations, not benchmarks.

## What we learned

The most important part of the harness was not the percentile function or the rate limiter. It was deciding what each counter means.

For messaging systems:

- connections are not messages;
- attempted publishes are not acknowledged publishes;
- acknowledged publishes are not recipient deliveries;
- group deliveries are not sender throughput;
- local smoke results are not production capacity.

Once those boundaries are explicit, optimization becomes easier too. A high connection P99, high ACK P99, and high group-delivery P99 point to different parts of the system.

The harness is open source in the [JuggleIM repository](https://github.com/juggleim/im-server?utm_source=devto&utm_medium=article&utm_campaign=juggleim_benchmark). If you work on real-time systems, feedback on the workload model, timing semantics, and missing scenarios is especially welcome.
