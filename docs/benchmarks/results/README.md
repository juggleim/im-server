# Verified local smoke run

These files prove that the complete benchmark path works against the repository's Docker Compose stack. They are development smoke results, not production capacity claims: the server and load generator shared one Apple M2 host, the measurement lasted only 10 seconds, and the working tree contained the new benchmark harness.

Run date: 2026-07-20 UTC. Server base commit: `b1f4ff52c3e5e6da2184f6b8ccb441fe04f13e21`. Database: the repository's default MySQL 8.0 Compose service. Workload: 20 WebSocket clients, stored and counted 256-byte messages, 3-second warm-up, 10-second measurement, 2-second delivery grace, target 50 publishes/second.

| Scenario | Connections | Successful ACKs | ACK rate | ACK P95 / P99 | Observed deliveries | Delivery rate | Delivery P95 / P99 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Private ring | 20/20 | 499/499 | 49.9/s | 2.37 / 17.09 ms | 499 | 49.9/s | 8.47 / 36.75 ms |
| 20-member group, 2 senders | 20/20 | 499/499 | 49.9/s | 4.54 / 30.19 ms | 9,492 | 949.2/s | 35.29 / 115.88 ms |

Machine-readable reports:

- [`local-smoke-20260720-private.json`](./local-smoke-20260720-private.json)
- [`local-smoke-20260720-group.json`](./local-smoke-20260720-group.json)

For publishable capacity results, follow the clean-commit, separate-load-generator, multi-trial, resource-monitoring, and full-environment checklist in [`BENCHMARKS.md`](../../../BENCHMARKS.md).
