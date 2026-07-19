# Launcher 配置参考

Launcher 默认从 `conf/config.yml` 读取 YAML。在 `launcher/` 目录启动时，
可通过以下参数指定其他文件：

```bash
go run main.go -config /path/to/config.yml
```

支持的配置项以 `commons/configures.ImConfig` 为准。未知 YAML 配置项会被忽略。

## 完整示例

以下示例包含当前支持的全部配置项。在本地开发环境之外使用前，请替换账号、
地址和密钥。

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

## 核心配置

| 配置项 | 用途 | 默认值或回退值 |
| --- | --- | --- |
| `nodeName` | 唯一的集群节点名称，也用于标识性能指标记录。 | 为空时生成短 UUID。集群部署应设置稳定且唯一的值。 |
| `nodeHost` | 向内部集群运行时声明的主机/IP，不是监听绑定地址。 | 为空时为 `127.0.0.1`。 |
| `msgStoreEngine` | 消息、历史消息和推送数据的存储实现：`mysql` 或 `mongo`。 | 为空时为 `mysql`。两种模式都会初始化 MySQL，用于共享业务数据。 |

## `log`

| 配置项 | 用途 | 默认值或回退值 |
| --- | --- | --- |
| `log.logPath` | 按小时轮转的应用日志目录。 | 代码未定义默认值；仓库示例使用 `./logs`。 |
| `log.logName` | 应用日志文件名前缀。 | 代码未定义默认值；仓库示例使用 `jim-info`。 |
| `log.logExpireHours` | 轮转应用日志的最长保留时间。 | 为零或负数时使用 `24`。 |

当前日志文件链接为 `<logPath>/<logName>.log`，轮转文件名中包含小时。

## `kvdb`

| 配置项 | 用途 | 默认值或回退值 |
| --- | --- | --- |
| `kvdb.isOpen` | 启用日志管理器使用的本地 LevelDB。 | 未应用显式回退值，因此默认为 `false`；仓库示例设置为 `true`。 |
| `kvdb.dataPath` | LevelDB 数据目录。 | 为空时使用 `<log.logPath>/kvdb_data`。 |

应使用节点本地的持久化路径。多个进程不能同时打开同一个 LevelDB 目录。

## `msglogs`

| 配置项 | 用途 | 默认值或回退值 |
| --- | --- | --- |
| `msglogs.logPath` | 每个应用消息日志的基础目录。 | 为空时使用 `<log.logPath>/msglogs/<appkey>`；设置后实际目录为 `<配置值>/msglogs/<appkey>`。 |
| `msglogs.maxBackups` | 保留的每小时消息日志备份数量。 | 为零或负数时使用 `24`。 |
| `msglogs.isCompress` | 是否压缩轮转后的消息日志。 | `false`。 |

## `mysql`

即使 `msgStoreEngine` 设置为 `mongo`，服务启动时仍然必须连接 MySQL。

| 配置项 | 用途 | 默认值或回退值 |
| --- | --- | --- |
| `mysql.user` | MySQL 用户名。 | 无。 |
| `mysql.password` | MySQL 密码。 | 无。 |
| `mysql.address` | `host:port` 格式的 MySQL 地址。 | 无。 |
| `mysql.name` | 数据库/Schema 名称。 | 无。 |
| `mysql.debug` | 启用 GORM SQL 日志。 | `false`。 |

Launcher 连接后会执行数据库升级。SQL 日志可能包含敏感业务数据，只应临时启用
`debug`。

## `mongodb`

仅当 `msgStoreEngine: mongo` 时初始化 MongoDB。

| 配置项 | 用途 | 默认值或回退值 |
| --- | --- | --- |
| `mongodb.address` | 不包含 `mongodb://` 前缀的 MongoDB 地址和可选凭据，例如 `127.0.0.1:27017` 或 `user:password@mongo:27017`。 | 无。 |
| `mongodb.name` | 消息集合使用的 MongoDB 数据库名称。 | 无。 |

Launcher 会在 `mongodb.address` 前添加 `mongodb://`。凭据中的保留字符需要进行
百分号编码。

## 监听服务和端口

| 配置项 | 服务 | 默认值或回退值 | 说明 |
| --- | --- | --- | --- |
| `apiGateway.httpPort` | 服务端 API | 为零或负数时使用 `9001`。 | 供可信业务服务器调用。 |
| `navGateway.httpPort` | 导航服务 | 无；为零时禁用独立 HTTP 监听。 | 向客户端 SDK 返回 WebSocket 连接地址。 |
| `connectManager.wsPort` | WebSocket | 为零或负数时使用 `9003`。 | 客户端 SDK 长连接。 |
| `adminGateway.httpPort` | 管理后台 | 为零或负数时使用 `8090`。 | 提供管理 API 和 Web 控制台。 |

这些服务都监听所有网卡（`:<port>`），`nodeHost` 不会限制绑定地址。Launcher
还会在不可配置的 `6060` 端口启动 Go `pprof` 监听。

## 管理和性能指标

| 配置项 | 用途 | 默认值或回退值 |
| --- | --- | --- |
| `adminSecret` | 作为服务端管理密钥传给内嵌管理后台。 | 无。它不是管理后台登录密码。 |
| `performanceMetrics.isOpen` | 每分钟保存一次节点运行指标，并保留 24 小时。 | `false`。 |

## Docker 和生产模板

镜像入口脚本将 `/opt/config_template.yaml` 渲染为 `/opt/conf/config.yml`。
应设置下列所有环境变量；缺少变量时 `run.sh` 只会输出警告，未替换的占位符仍会
保留在生成的 YAML 中。

| 环境变量 | 生成的 YAML 配置项 | 本地 Docker Compose 值 |
| --- | --- | --- |
| `POD_NAME` | `nodeName` | `jim-node-1` |
| `POD_IP` | `nodeHost` | `127.0.0.1` |
| `MSG_STORE_ENGINE` | `msgStoreEngine` | `mysql` |
| `MYSQL_ROOT_PASSWORD` | `mysql.password` | `${MYSQL_ROOT_PASSWORD:-juggleim}` |
| `MYSQL_ADDR` | `mysql.address` | `mysql:3306` |
| `MYSQL_DB_NAME` | `mysql.name` | `jim_db` |
| `MONGODB_ADDR` | `mongodb.address` 的主机部分 | `mongo`（MySQL 模式不使用） |
| `MONGODB_ROOT_PASSWORD` | `mongodb.address` 的密码部分 | `unused`（MySQL 模式不使用） |
| `WS_PORT` | `connectManager.wsPort` | `9003` |
| `API_HTTP_PORT` | `apiGateway.httpPort` | `9001` |
| `NAV_HTTP_PORT` | `navGateway.httpPort` | `9002` |
| `ADMIN_HTTP_PORT` | `adminGateway.httpPort` | `8090` |

不同部署方式的差异：

| 模式 | 配置来源 | 重要行为 |
| --- | --- | --- |
| 本地 Docker Compose | `docker-compose.yml` 环境变量经 `launcher/scripts/config_template.yaml` 渲染 | 使用 `mysql` 服务名，在 `9002` 启用导航服务，并发布 `9001`、`9002`、`9003`、`8090` 和 `6060`。 |
| 源码启动 | `launcher/conf/config.yml`，除非传入 `-config` | 使用回环数据库地址。仓库中的文件未配置 `navGateway`，因此添加配置前独立导航监听处于禁用状态。 |
| 生产镜像 | 部署系统提供的环境变量由 `launcher/scripts/run.sh` 渲染 | 需要提供唯一节点标识、可达服务地址、密钥、持久化存储和明确的网络访问控制。 |

`launcher/conf/config.yml` 中的 `defaultPort` 和容器模板中的 `imApiDomain`
不属于 Launcher 当前的 `ImConfig`。Launcher 监听端口应通过上述 Gateway 配置，
内嵌 API 域名由 `apiGateway.httpPort` 派生。

## 安全检查清单

- 不要提交真实的 MySQL/MongoDB 密码、`adminSecret`、租户 App Secret 或生成的
  配置文件。
- 平台支持时，应使用密钥管理服务或权限受限的只读挂载配置，避免使用明文环境
  变量。
- 不要将管理后台（`8090`）或 `pprof`（`6060`）暴露到公网；服务端 API
  （`9001`）应只允许可信业务服务访问。
- 公网导航服务和 WebSocket 应部署在 TLS、防火墙或可信反向代理之后；Launcher
  本身提供的是明文 HTTP/WS 监听。
- 任何非本地部署都必须替换本地管理后台账号密码（`admin` / `123456`），并保护
  通过管理 API 返回的租户 App Key/App Secret。
