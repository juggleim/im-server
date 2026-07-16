# Goroutine configuration and saturation verification

## Configuration inventory

| Component | Before | After / behavior |
| --- | ---: | --- |
| Actor callback pool | 8,192 workers | 256 workers by default; `IM_ACTOR_CALLBACK_WORKERS` overrides it |
| Actor shared executor pool | 16,384 workers | 1,024 workers by default; `IM_ACTOR_EXECUTOR_WORKERS` overrides it |
| Actor admission queue | 8,192 items per executor | unchanged; now provides real backpressure |
| Standalone message actors | 3,072 / 6,144 workers per configured actor | unchanged; explicitly configured by each service |
| Message dispatch shards | up to 8,192 lazy workers | unchanged; bounded and preserves per-user ordering |
| Generic async task pool | 64 workers + 10,000 queued items | unchanged; bounded submission |
| Conversation batch executors | up to 128 lazy ticker workers | unchanged while running; `Stop` now terminates them |

Invalid, empty, or non-positive Actor worker environment values fall back to
the defaults above.

## Root causes and changes

The Actor dispatcher used `go pool.Process(...)`. Since `Process` blocks when
all workers are occupied, the extra `go` bypassed the 8,192-item admission
queue and created one blocked goroutine for every excess request. Actor pool
calls now use a shared slot semaphore: active submitters cannot exceed pool
concurrency, queued work remains bounded, and configured parallelism is kept.

Queue capacity and worker concurrency are now separate. The old defaults
prestarted 24,576 workers per Actor system; the new defaults prestart 1,280.

The batch executor previously ranged over `Ticker.C`. `Ticker.Stop()` does not
close that channel, so stopped executors leaked their goroutine. It now selects
on an explicit done channel and waits for termination during `Stop`.

## Directed saturation results

```bash
go test ./commons/gmicro/actorsystem ./commons/tools \
  -run 'Test(Actor|BatchExecutor)' -count=1 -v
go test -race ./commons/gmicro/actorsystem ./commons/tools \
  -run 'Test(Actor|BatchExecutor)' -count=1
```

Observed on 2026-07-16:

| Scenario | Before | After |
| --- | ---: | ---: |
| 1 blocked worker + burst of 2,000 Actor requests | +2,000 goroutines | +0 goroutines beyond the saturated baseline |
| Default Actor dispatcher startup | 24,576 configured workers | 1,280 configured workers; 1,282 observed goroutines including control loops |
| 8-worker execution, 80 jobs at 5 ms/job | not used as a failure baseline | all 8 workers reached concurrently; about 59 ms total |
| Stop 64 batch executors | ticker goroutines remained blocked | returned to baseline (within test-process tolerance) |

`go test ./...` also exercises integration suites that require the repository's
external services and credentials. Core Actor, tools, and affected message
packages can be verified independently with the commands above and their
package-level test commands.
