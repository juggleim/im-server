# Docker Compose 故障排查

本文适用于仓库根目录 [`docker-compose.yml`](../docker-compose.yml) 定义的本地环境。所有命令均在仓库根目录执行。

## 1. 检查容器状态

```bash
docker compose ps
```

`mysql` 服务应先变为 `healthy`，随后 `im-server` 才会启动。如果服务缺失或已经退出，显示包括已停止容器在内的完整状态：

```bash
docker compose ps --all
```

如果当前 Docker Compose 支持 `--wait`，可启动或重建服务并等待健康检查：

```bash
docker compose up -d --wait --wait-timeout 120
```

如果不支持 `--wait`，先执行 `docker compose up -d`，然后重复执行 `docker compose ps`，直到 MySQL 为健康状态且 `im-server` 正在运行。

## 2. 查看日志

查看 MySQL 和服务端最近的日志：

```bash
docker compose logs --tail=200 mysql
docker compose logs --tail=200 im-server
```

复现问题时持续查看两个服务的日志：

```bash
docker compose logs --follow mysql im-server
```

MySQL 常见问题包括首次导入表结构未完成、数据目录损坏，或数据卷创建时使用的密码与当前密码不同。`im-server` 常见问题包括仍在等待 MySQL、生成的配置无效，或宿主机端口已被占用。

## 3. 排查 MySQL 健康检查

查看 Docker 记录的健康检查输出：

```bash
docker inspect --format '{{range .State.Health.Log}}{{println .End .ExitCode .Output}}{{end}}' juggleim-mysql
```

在容器内部确认 MySQL 可以接受连接，同时避免将密码写入命令历史：

```bash
docker compose exec mysql mysqladmin ping --host=localhost --user=root --password
```

根据提示输入通过 `MYSQL_ROOT_PASSWORD` 配置的本地密码。检查成功时会显示 `mysqld is alive`。

首次启动时，MySQL 会导入 [`sql/imserver.sql`](../sql/imserver.sql) 和 [`docker/mysql-init`](../docker/mysql-init/) 下的脚本。在冷启动环境中初始化可能需要更长时间。不要因为最初几次健康检查仍在重试就直接删除数据卷，应先查看日志。

## 4. 解决端口冲突

本地环境会映射以下宿主机端口：

| 端口 | 服务 |
| ---: | :--- |
| `3306` | MySQL |
| `9001` | 服务端 API |
| `9002` | 导航服务 |
| `9003` | WebSocket |
| `8090` | 管理后台 |

在 macOS 或 Linux 上，可通过以下命令查找监听进程：

```bash
for port in 3306 9001 9002 9003 8090; do
  echo "Port ${port}"
  lsof -nP -iTCP:"${port}" -sTCP:LISTEN
done
```

停止冲突的本地进程或容器。也可以只修改 `docker-compose.yml` 端口映射的宿主机一侧。例如使用 `"13306:3306"` 将 MySQL 发布到宿主机的 `13306`，同时保留容器端口 `3306`。服务之间通过 Compose 网络通信，内部地址应继续使用 `mysql:3306`。

## 5. 确认本地端口可访问

确认管理后台可以返回 HTTP 响应：

```bash
curl --fail --silent --show-error --output /dev/null http://127.0.0.1:8090/
```

确认 WebSocket 及其他映射端口能够建立 TCP 连接：

```bash
for port in 9001 9002 9003 8090; do
  nc -vz 127.0.0.1 "${port}"
done
```

TCP 检查只能确认端口可达。成功建立 WebSocket 会话仍需要有效的应用、用户 Token 和客户端握手。

如果服务可以通过 `127.0.0.1` 访问，但无法从其他机器访问，请检查宿主机防火墙、云安全组、反向代理和已配置的公网连接地址。不要将 MySQL 或管理后台直接暴露到公网。

## 6. 安全停止或重置环境

停止并删除容器和 Compose 网络，同时保留 MySQL 数据：

```bash
docker compose down
```

停止服务并永久删除 MySQL 数据卷：

```bash
docker compose down -v
```

仅在确实需要全新本地数据库时使用 `down -v`。该命令会删除 Compose 数据卷中的本地应用、用户、消息、配置及其他全部记录，请先备份需要保留的数据。

重置后重新构建并启动：

```bash
docker compose up -d --build
```

如果问题仍然存在，请提交 GitHub Issue，并附上操作系统、Docker 与 Compose 版本、`docker compose ps --all` 输出以及相关的脱敏日志。发布前应删除密码、应用密钥、Token、公网 IP、用户数据和消息内容。
