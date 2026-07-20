# Launcher configuration reference

The launcher reads YAML from `conf/config.yml` by default. When starting from
the `launcher/` directory, use another file with:

```bash
go run main.go -config /path/to/config.yml
```

The source of truth for supported keys is `commons/configures.ImConfig`.
Unknown YAML keys are ignored.

## Complete example

This example includes every currently supported key. Replace credentials,
addresses, and secrets before using it outside a local development machine.

```yaml
nodeName: node-1
nodeHost: 127.0.0.1
msgStoreEngine: mysql

log:
  logPath: ./logs
  logName: jim-info
  logExpireHours: 24

kvdb:
  isOpen: true
  dataPath: ./logs/kvdb_data

msglogs:
  logPath: ./logs
  maxBackups: 24
  isCompress: false

mysql:
  user: root
  password: change-me
  address: 127.0.0.1:3306
  name: jim_db
  debug: false

mongodb:
  address: 127.0.0.1:27017
  name: jim_msgs

connectManager:
  wsPort: 9003

apiGateway:
  httpPort: 9001

navGateway:
  httpPort: 9002

adminGateway:
  httpPort: 8090

adminSecret: change-me

performanceMetrics:
  isOpen: false
```

## Core settings

| Key | Purpose | Default or fallback |
| --- | --- | --- |
| `nodeName` | Unique cluster node name. It also identifies performance-metric rows. | A generated short UUID when empty. Set a stable, unique value in clustered deployments. |
| `nodeHost` | Host/IP advertised to the internal cluster runtime. This is not a listener bind address. | `127.0.0.1` when empty. |
| `msgStoreEngine` | Message, history, and push storage implementation: `mysql` or `mongo`. | `mysql` when empty. MySQL is still initialized for shared application data in both modes. |

## `log`

| Key | Purpose | Default or fallback |
| --- | --- | --- |
| `log.logPath` | Directory for hourly application logs. | No code-defined default. The shipped files use `./logs`. |
| `log.logName` | Application log filename prefix. | No code-defined default. The shipped files use `jim-info`. |
| `log.logExpireHours` | Maximum age of rotated application logs. | `24` when zero or negative. |

The active file is linked as `<logPath>/<logName>.log`; rotated files include
the hour in their names.

## `kvdb`

| Key | Purpose | Default or fallback |
| --- | --- | --- |
| `kvdb.isOpen` | Enables the local LevelDB store used by the log manager. | `false` because no explicit fallback is applied. The shipped files set it to `true`. |
| `kvdb.dataPath` | LevelDB directory. | `<log.logPath>/kvdb_data` when empty. |

Use a node-local, persistent path. Multiple processes must not open the same
LevelDB directory.

## `msglogs`

| Key | Purpose | Default or fallback |
| --- | --- | --- |
| `msglogs.logPath` | Base directory for per-application message logs. | `<log.logPath>/msglogs/<appkey>` when empty. When set, the effective directory is `<value>/msglogs/<appkey>`. |
| `msglogs.maxBackups` | Number of hourly message-log backups retained. | `24` when zero or negative. |
| `msglogs.isCompress` | Compresses rotated message logs. | `false`. |

## `mysql`

MySQL is required at startup even when `msgStoreEngine` is `mongo`.

| Key | Purpose | Default or fallback |
| --- | --- | --- |
| `mysql.user` | MySQL user. | None. |
| `mysql.password` | MySQL password. | None. |
| `mysql.address` | MySQL endpoint in `host:port` form. | None. |
| `mysql.name` | Database/schema name. | None. |
| `mysql.debug` | Enables GORM SQL logging. | `false`. |

The launcher runs database upgrades after connecting. Enable `debug` only
temporarily because SQL logs can contain sensitive application data.

## `mongodb`

MongoDB is initialized only when `msgStoreEngine: mongo`.

| Key | Purpose | Default or fallback |
| --- | --- | --- |
| `mongodb.address` | MongoDB authority and optional credentials, without the `mongodb://` prefix. Examples: `127.0.0.1:27017` or `user:password@mongo:27017`. | None. |
| `mongodb.name` | MongoDB database name for message collections. | None. |

The launcher prepends `mongodb://` to `mongodb.address`. Percent-encode
reserved characters in credentials.

## Listeners and ports

| Key | Service | Default or fallback | Notes |
| --- | --- | --- | --- |
| `apiGateway.httpPort` | Server API | `9001` when zero or negative. | Called by trusted business servers. |
| `navGateway.httpPort` | Navigator | None; zero disables its separate HTTP listener. | Returns WebSocket connection addresses to client SDKs. |
| `connectManager.wsPort` | WebSocket | `9003` when zero or negative. | Long-lived client SDK connections. |
| `adminGateway.httpPort` | Admin | `8090` when zero or negative. | Hosts the admin API and web console. |

These services listen on all interfaces (`:<port>`). `nodeHost` does not limit
the bind address. The launcher also starts an unconfigured Go `pprof` listener
on port `6060`.

## Admin and metrics

| Key | Purpose | Default or fallback |
| --- | --- | --- |
| `adminSecret` | Passed to the embedded admin console as its server-side admin secret. | None. This is not the admin login password. |
| `performanceMetrics.isOpen` | Persists per-node runtime metrics once per minute and retains 24 hours. | `false`. |

## Docker and production template

The image entrypoint renders `/opt/config_template.yaml` to
`/opt/conf/config.yml`. Set every variable below; `run.sh` warns about missing
values but leaves unresolved placeholders in the generated YAML.

| Environment variable | Generated YAML key | Local Docker Compose value |
| --- | --- | --- |
| `POD_NAME` | `nodeName` | `jim-node-1` |
| `POD_IP` | `nodeHost` | `127.0.0.1` |
| `MSG_STORE_ENGINE` | `msgStoreEngine` | `mysql` |
| `MYSQL_ROOT_PASSWORD` | `mysql.password` | `${MYSQL_ROOT_PASSWORD:-juggleim}` |
| `MYSQL_ADDR` | `mysql.address` | `mysql:3306` |
| `MYSQL_DB_NAME` | `mysql.name` | `jim_db` |
| `MONGODB_ADDR` | Host portion of `mongodb.address` | `mongo` (unused in MySQL mode) |
| `MONGODB_ROOT_PASSWORD` | Password portion of `mongodb.address` | `unused` (unused in MySQL mode) |
| `WS_PORT` | `connectManager.wsPort` | `9003` |
| `API_HTTP_PORT` | `apiGateway.httpPort` | `9001` |
| `NAV_HTTP_PORT` | `navGateway.httpPort` | `9002` |
| `ADMIN_HTTP_PORT` | `adminGateway.httpPort` | `8090` |

Deployment differences:

| Mode | Configuration source | Important behavior |
| --- | --- | --- |
| Local Docker Compose | `docker-compose.yml` environment values rendered through `launcher/scripts/config_template.yaml` | Uses the `mysql` service name, enables Navigator on `9002`, publishes `9001`, `9002`, `9003`, `8090`, and `6060`. |
| From source | `launcher/conf/config.yml`, unless `-config` is provided | Uses loopback database addresses. The checked-in file omits `navGateway`, so the separate Navigator listener is disabled until configured. |
| Production image | Deployment-supplied environment values rendered by `launcher/scripts/run.sh` | Supply unique node identity, reachable service addresses, secrets, persistent storage, and explicit network controls. |

`defaultPort` in `launcher/conf/config.yml` and `imApiDomain` in the container
template are not fields in the launcher's current `ImConfig`. Configure
launcher listeners through the gateway sections above; the launcher derives
the embedded API domain from `apiGateway.httpPort`.

## Security checklist

- Never commit real MySQL/MongoDB passwords, `adminSecret`, tenant app secrets,
  or generated configuration files.
- Use a secret manager or read-only mounted configuration with restrictive file
  permissions instead of plain environment variables where your platform
  supports it.
- Do not expose the Admin listener (`8090`) or `pprof` (`6060`) to the public
  Internet. Restrict the Server API (`9001`) to trusted business services.
- Put public Navigator and WebSocket endpoints behind TLS, a firewall, or a
  trusted reverse proxy. The launcher listeners themselves are plain HTTP/WS.
- Replace the local admin console credentials (`admin` / `123456`) before any
  non-local deployment and protect tenant app keys/secrets returned through the
  admin API.
