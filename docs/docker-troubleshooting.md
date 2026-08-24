# Docker Compose troubleshooting

This guide covers the local stack defined by [`docker-compose.yml`](../docker-compose.yml). Run all commands from the repository root.

## 1. Check container status

```bash
docker compose ps
```

The `mysql` service should become `healthy` before `im-server` starts. If either service is missing or stopped, include stopped containers in the output:

```bash
docker compose ps --all
```

Start or recreate the stack and wait for its health checks when your Docker Compose version supports `--wait`:

```bash
docker compose up -d --wait --wait-timeout 120
```

If `--wait` is unavailable, use `docker compose up -d`, then rerun `docker compose ps` until MySQL is healthy and `im-server` is running.

## 2. Inspect logs

Show the latest MySQL and server logs:

```bash
docker compose logs --tail=200 mysql
docker compose logs --tail=200 im-server
```

Follow both services while reproducing a failure:

```bash
docker compose logs --follow mysql im-server
```

Common MySQL failures include an incomplete first-time schema import, a damaged data directory, or a password that differs from the value used when the volume was created. Common `im-server` failures include waiting for MySQL, invalid generated configuration, or a host port already in use.

## 3. Diagnose the MySQL health check

Inspect the health-check output recorded by Docker:

```bash
docker inspect --format '{{range .State.Health.Log}}{{println .End .ExitCode .Output}}{{end}}' juggleim-mysql
```

Confirm that MySQL accepts connections inside its container without printing a password in the shell history:

```bash
docker compose exec mysql mysqladmin ping --host=localhost --user=root --password
```

Enter the local password configured through `MYSQL_ROOT_PASSWORD`. A successful check prints `mysqld is alive`.

On first startup, MySQL imports [`sql/imserver.sql`](../sql/imserver.sql) and the scripts under [`docker/mysql-init`](../docker/mysql-init/). Initialization can take longer on a cold machine. Do not delete the volume merely because the first health checks are still retrying; inspect the logs first.

## 4. Resolve port conflicts

The local stack publishes these host ports:

| Port | Service |
| ---: | :--- |
| `3306` | MySQL |
| `9001` | Server API |
| `9002` | Navigator |
| `9003` | WebSocket |
| `8090` | Admin console |

On macOS or Linux, find listeners with:

```bash
for port in 3306 9001 9002 9003 8090; do
  echo "Port ${port}"
  lsof -nP -iTCP:"${port}" -sTCP:LISTEN
done
```

Stop the conflicting local process or container. Alternatively, override only the conflicting host port. For example, publish MySQL on host port `13306` while keeping its container port `3306`:

```bash
MYSQL_HOST_PORT=13306 docker compose up -d
```

The available overrides are `MYSQL_HOST_PORT`, `API_HOST_PORT`, `NAV_HOST_PORT`, `WS_HOST_PORT`, `ADMIN_HOST_PORT`, and `PPROF_HOST_PORT`. The services communicate over the Compose network and continue to use their original container ports internally.

## 5. Confirm local reachability

Verify that the admin console returns an HTTP response:

```bash
curl --fail --silent --show-error --output /dev/null http://127.0.0.1:8090/
```

Verify that the WebSocket and other published TCP ports accept connections:

```bash
for port in 9001 9002 9003 8090; do
  nc -vz 127.0.0.1 "${port}"
done
```

This TCP check confirms reachability only. A successful WebSocket session still requires a valid application, user token, and client handshake.

If the services work on `127.0.0.1` but not from another machine, check the host firewall, cloud security-group rules, reverse proxy, and the configured public connection address. Do not expose MySQL or the admin console directly to the public Internet.

## 6. Stop or reset the stack safely

Stop and remove the containers and Compose network while preserving MySQL data:

```bash
docker compose down
```

Stop the stack and permanently remove its MySQL volume:

```bash
docker compose down -v
```

Use `down -v` only when you intentionally want a clean local database. It deletes local applications, users, messages, configuration, and every other record stored in the Compose volume. Back up any data you need first.

After a reset, rebuild and start the stack:

```bash
docker compose up -d --build
```

If the problem remains, open a GitHub issue and include the operating system, Docker and Compose versions, `docker compose ps --all`, and relevant sanitized logs. Remove passwords, application secrets, tokens, public IP addresses, and user or message data before posting.
