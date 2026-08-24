# Reproducible benchmarks

JuggleIM's benchmark harness measures connection establishment separately from steady-state message delivery. It favors reproducibility and explicit limitations over headline numbers.

The harness runs against an isolated JuggleIM deployment, creates synthetic users through the signed server API, connects real WebSocket clients using the production protocol, and records a machine-readable JSON report. It must never target a public or production deployment.

## Workloads

### Private chat

Every connected client is a sender. Clients form a ring: client `N` sends to client `N+1`, and the final client sends to the first. The configured rate is the aggregate publish rate across the ring.

### Group chat

All clients join one synthetic group. A configurable subset of clients publishes to the group at the aggregate target rate. The report distinguishes accepted publishes from observed recipient deliveries, making fan-out visible without treating it as sender throughput.

Both workloads use stored, counted text messages by default. Each JSON payload contains a run ID, phase, sequence number, nanosecond send timestamp, and padding to reach the configured minimum size.

## Metrics

The JSON report contains:

- connection attempts, successes, failures, error codes, and connection latency;
- publish-to-ACK throughput and P50/P95/P99 latency;
- publish-to-recipient-callback throughput and P50/P95/P99 latency;
- workload settings, including clients, senders, message rate, payload size, warm-up, duration, and delivery grace;
- server commit, dirty-tree state, OS, architecture, Go version, CPU model, logical CPU count, memory, database label, and deployment label;
- explicit limitations that must accompany any published result.

Warm-up traffic is excluded from all message measurements. Connection establishment is completed before warm-up and is never mixed into steady-state latency.

A small **[verified local smoke run](./docs/benchmarks/results/README.md)** and its machine-readable private-chat and group-chat reports are checked into the repository. They validate the full harness; they are explicitly not production capacity claims.

## One-command local run

Requirements:

- Docker with Compose;
- Go version from [`go.mod`](./go.mod);
- `curl`, `jq`, `openssl`, and `shasum`.

Run both baseline workloads:

```bash
scripts/run-benchmark.sh all
```

The script builds and starts the local Compose stack, creates an isolated application, runs private and group workloads, and writes JSON files under `benchmark-results/`. It leaves the stack running so logs, profiles, and database state can be inspected. Stop it without deleting the database volume:

```bash
docker compose down
```

The default baseline is intentionally modest:

| Setting | Default |
| --- | ---: |
| Connected clients | 50 |
| Aggregate publish rate | 200 messages/second |
| Group senders | 2 |
| Warm-up | 10 seconds |
| Measurement | 30 seconds |
| Delivery grace | 3 seconds |
| Minimum payload | 256 bytes |

Override settings with environment variables:

```bash
BENCH_CLIENTS=500 \
BENCH_RATE=1000 \
BENCH_WARMUP=30s \
BENCH_DURATION=2m \
BENCH_PAYLOAD_BYTES=1024 \
scripts/run-benchmark.sh all
```

Useful variables are `BENCH_CLIENTS`, `BENCH_RATE`, `BENCH_GROUP_SENDERS`, `BENCH_WARMUP`, `BENCH_DURATION`, `BENCH_DELIVERY_GRACE`, `BENCH_PAYLOAD_BYTES`, and `BENCH_OUTPUT_DIR`.

To avoid conflicting with local development services, the runner uses host ports `19001` (API), `19002` (navigator), `19003` (WebSocket), `18090` (admin), `16060` (pprof), and `13306` (MySQL). Override them with the matching `*_HOST_PORT` variables. Containers still use JuggleIM's standard ports over the internal Compose network.

## Running the CLI against an isolated environment

The CLI reads the application secret only from `JIM_BENCH_APP_SECRET`, keeping it out of command arguments and result files:

```bash
export JIM_BENCH_APP_KEY='benchmark-app'
export JIM_BENCH_APP_SECRET='replace-with-isolated-app-secret'

go run ./cmd/jimbench \
  --scenario private \
  --clients 100 \
  --rate 500 \
  --warmup 30s \
  --duration 2m \
  --payload-bytes 256
```

Loopback targets are enforced by default. `--allow-non-loopback` exists only for a dedicated, isolated benchmark environment. Never use it with shared staging, public, or production infrastructure.

## Reproduction checklist

Before publishing a result:

1. Use a clean, named server commit and include the generated JSON files.
2. Run the load generator on a separate host for serious capacity measurements.
3. Record server and database hardware separately in `--environment` and `--database` labels.
4. State container CPU/memory limits, MySQL configuration, network topology, TLS status, and storage type.
5. Run at least three trials after warm-up and publish every trial, not only the best one.
6. Report errors and saturation behavior alongside throughput.
7. Keep private-chat and group fan-out results separate.
8. Use server-side CPU, memory, network, disk, database, and pprof data to explain bottlenecks.

Example for a dedicated environment:

```bash
go run ./cmd/jimbench \
  --scenario group \
  --ws-url wss://benchmark-im.example.internal \
  --api-url https://benchmark-api.example.internal/apigateway \
  --allow-non-loopback \
  --clients 1000 \
  --group-senders 10 \
  --rate 1000 \
  --warmup 1m \
  --duration 5m \
  --environment 'server: 8 vCPU/16 GiB, load generator: 8 vCPU/16 GiB, 10 GbE, TLS enabled' \
  --database 'MySQL 8.0.42, 8 vCPU/32 GiB, local NVMe, configuration linked in report notes'
```

## Interpreting results

- ACK throughput is the rate at which JuggleIM accepts publishes; it is not the same as recipient deliveries.
- Delivery throughput counts callbacks observed by connected clients. Group fan-out can therefore be much larger than publish throughput.
- ACK latency includes client serialization, WebSocket transport, server processing, and the ACK return path.
- Delivery latency includes the full sender-to-recipient path but assumes synchronized timing because the load generator hosts all benchmark clients.
- A local Compose result is a development baseline, not a production capacity claim.

Do not compare JuggleIM with another project unless both systems use equivalent delivery semantics, persistence, acknowledgements, payloads, fan-out, hardware, databases, and test duration.

## Development checks

```bash
go test ./simulator/benchmark ./cmd/jimbench
go test -race ./simulator/benchmark
go vet ./simulator/benchmark ./cmd/jimbench
```
