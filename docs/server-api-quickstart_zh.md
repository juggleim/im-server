# 服务端 API 快速开始

本文从一个全新的本地 Docker Compose 环境开始，完成客户端 SDK 建立连接前所需的最小服务端流程：

1. 登录本地管理 API。
2. 创建应用（租户）。
3. 注册两个合成测试用户并获取客户端 Token。
4. 通过服务端 API 发送一条私聊消息。

该示例已经在当前 [`docker-compose.yml`](../docker-compose.yml) 上完成验证，仅使用自动生成的本地标识和测试消息内容。

## 前置条件

安装 Docker Compose、`curl`、`jq` 和 `openssl`，然后启动本地环境：

```bash
docker compose up -d --wait --wait-timeout 120
docker compose ps
```

MySQL 应为 `healthy`，`im-server` 应处于运行状态。如果没有正常启动，请参考 [Docker Compose 故障排查指南](./docker-troubleshooting_zh.md)。

## 运行已验证的示例

```bash
bash examples/server-api-quickstart.sh
```

脚本使用 `docker-compose.yml` 中仅供本地开发的默认管理账号，创建一个唯一应用，注册 Alice 和 Bob，并由 Alice 向 Bob 发送一条持久化的 `jg:text` 私聊消息。脚本会输出生成的用户 ID、用户 Token 和消息 ID，但不会输出应用密钥。

无需修改脚本即可覆盖本地地址和管理账号：

```bash
ADMIN_BASE_URL=http://127.0.0.1:8090/admingateway \
API_BASE_URL=http://127.0.0.1:9001/apigateway \
ADMIN_ACCOUNT=admin \
ADMIN_PASSWORD='<本地管理密码>' \
bash examples/server-api-quickstart.sh
```

## 使用的 API

| 步骤 | 接口 | 结果 |
| :--- | :--- | :--- |
| 管理后台登录 | `POST /admingateway/login` | 短期有效的本地管理授权 |
| 创建应用 | `POST /admingateway/apps/create` | 应用 Key 和仅限服务端保存的应用密钥 |
| [注册用户](https://www.juggle.im/docs/server/user/register/) | `POST /apigateway/users/register` | 用户 ID 和客户端 Token |
| [发送私聊消息](https://www.juggle.im/docs/server/message/privatemsg/) | `POST /apigateway/messages/private/send` | 每个接收者对应的消息 ID |

服务端 API 请求使用文档定义的[签名请求头](https://www.juggle.im/docs/server/api/#header)：`appkey`、`nonce`、`timestamp` 和 `signature`。签名是以下内容的 SHA-1 摘要：

```text
app_secret + nonce + timestamp
```

示例会为每个请求重新计算签名。真实接入中，这些操作由业务服务器完成，业务服务器只向客户端返回其所需的用户 Token 和连接信息。

## 连接客户端 SDK

同时使用以下三个值：

- 生成的应用 Key；
- 一个用户 ID 及其匹配的 Token；
- 本地 WebSocket 地址 `ws://127.0.0.1:9003`。

根据目标平台的[客户端 SDK 快速开始](https://www.juggle.im/docs/client/quickstart/android/)完成初始化。不要将其他应用或其他用户的 Token 与当前应用 Key、用户 ID 混用。

## 安全边界

以下值和操作必须保留在可信服务端：

- 管理账号密码和管理授权；
- 应用密钥；
- 服务端 API 签名逻辑；
- 创建应用、注册用户或发送高权限服务端消息的调用。

客户端可以获取自身所需的应用 Key、用户 ID、短期用户 Token，以及导航服务/WebSocket 地址。禁止将应用密钥或管理授权写入 Web、桌面或移动客户端代码，也不要记录到日志或提交到版本库。

脚本会输出合成测试用户的本地 Token，方便复制到本地 SDK 测试环境。生产用户 Token 应按凭据管理，并且只能交付给它对应的已认证用户。
