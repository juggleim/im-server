<div align="center">

<img height="120" src="./docs/logo.png" alt="JuggleIM Logo">

# JuggleIM

**A high-performance, scalable open-source instant messaging (IM) system.**

[![License](https://img.shields.io/github/license/juggleim/im-server?color=yellow&style=flat-square)](./LICENSE)
[![Go](https://img.shields.io/badge/go-%3E%3D1.20-30dff3?style=flat-square&logo=go)](https://github.com/juggleim/im-server)
[![Release](https://img.shields.io/github/v/release/juggleim/im-server?style=flat-square&color=brightgreen)](https://github.com/juggleim/im-server/releases)
[![Stars](https://img.shields.io/github/stars/juggleim/im-server?style=flat-square&color=orange)](https://github.com/juggleim/im-server/stargazers)
[![Forks](https://img.shields.io/github/forks/juggleim/im-server?style=flat-square)](https://github.com/juggleim/im-server/network/members)
[![Last Commit](https://img.shields.io/github/last-commit/juggleim/im-server?style=flat-square)](https://github.com/juggleim/im-server/commits)

**[简体中文](./README_zh.md)** | **English**

**[Website](https://www.juggle.im)** ·
**[Docs](https://www.juggle.im/docs/guide/intro/)** ·
**[Quick Deploy](https://www.juggle.im/docs/guide/deploy/quickdeploy/)** ·
**[API Reference](https://www.juggle.im/docs/server/api/)** ·
**[Ask a Question](https://github.com/juggleim/im-server/issues)**

If this project helps you, please give it a ⭐ **Star** — it keeps the project easy to find and motivates us to keep building!

</div>

---

## 📖 What is JuggleIM

JuggleIM is a **ready-to-use, self-hostable** instant messaging (IM) backend. Built on a Protobuf + WebSocket long-connection protocol, it focuses on efficient message delivery and reliable storage, letting you add chat capabilities to your app, website, or business system in minutes.

Whether you are building a social product, a customer-service system, IoT device communication, live-stream chat, or AI bot conversations, JuggleIM serves as a solid messaging foundation. It is **multi-tenant** by design — a single deployment can host multiple fully isolated applications — and the professional edition scales horizontally to support **hundreds of millions of daily active users**.

> Want to try it right away? Check the [official docs](https://www.juggle.im/docs/guide/intro/) or jump to [Quick Deploy](#-quick-deploy) below.

## ✨ Key Features

**🚀 High Performance & High Availability**
- Protobuf + WebSocket long connections — low bandwidth, high throughput, and reliable connectivity even on poor networks
- Professional edition supports clustered deployment with unlimited horizontal scaling, powering apps with hundreds of millions of DAU
- Handles large groups of 10,000–100,000 members without losing messages, plus unlimited-size live chat rooms

**🔒 Secure & Reliable**
- End-to-end encryption of protocol and data — no leakage risk
- Multi-device online presence and message sync keep state consistent across all endpoints

**🌍 Flexible Deployment & Global Reach**
- Supports public cloud, private cloud, and managed cloud deployment models
- Global link acceleration for worldwide-scale applications

**🧩 Easy to Integrate & Extend**
- SDKs for Android, iOS, Web, PC, Flutter, and HarmonyOS, each with a demo and docs
- Rich REST APIs and WebHooks for integrating with your existing systems
- Built-in AI bot connectivity — easily plug in large language models
- Comes with ops tooling and an admin console for simple maintenance

## 🗂 Table of Contents

- [Ecosystem](#-ecosystem)
- [Quick Deploy](#-quick-deploy)
- [Manual Deployment](#-manual-deployment)
- [Create an Application (Tenant)](#-create-an-application-tenant)
- [Integration](#-business-server--client-integration)
- [Community](#-community)
- [Star History](#-star-history)

## 🧬 Ecosystem

JuggleIM follows a layered architecture — "core service + business service + multi-platform SDKs + demos" — with clearly separated repositories that can be composed and customized as needed.

| Repository | Description |
| :--- | :--- |
| **[im-server](https://github.com/juggleim/im-server/)** | Core IM service handling message delivery, storage, and related IM logic (this repo) |
| [jugglechat-server](https://github.com/juggleim/jugglechat-server) | Demo business service handling user registration/login, group creation, friends, etc. — a base for your own development |
| [jugglechat-server-java](https://github.com/juggleim/jugglechat-server-java) | Java version of the demo business service |
| [imserver-console](https://github.com/juggleim/imserver-console) | Admin console for IM configuration and business metrics |
| [imsdk-android](https://github.com/juggleim/imsdk-android) | Android imsdk with a UI demo, ready for customization |
| [imsdk-ios](https://github.com/juggleim/imsdk-ios) | iOS imsdk with a UI demo, ready for customization |
| [imsdk-web](https://github.com/juggleim/imsdk-web) | Web imsdk |
| [imsdk-flutter](https://github.com/juggleim/imsdk-flutter) | Flutter version of imsdk |
| [imsdk-harmony](https://github.com/juggleim/imsdk-harmony) | HarmonyOS imsdk with a UI demo, ready for customization |
| [jugglechat-web](https://github.com/juggleim/jugglechat-web) | Web demo integrating imsdk-web, ready for customization |
| [jugglechat-desktop](https://github.com/juggleim/jugglechat-desktop) | Desktop demo integrating imsdk-pc, ready for customization |
| [jugglelive-web](https://github.com/juggleim/jugglelive-web) | Live chat-room demo integrating imsdk-web, ready for customization |
| [bot-connector](https://github.com/juggleim/bot-connector) | Bot connector service bridging im-server and third-party bots |
| [imserver-sdk-go](https://github.com/juggleim/imserver-sdk-go) | SDK wrapping the im-server server-side API for easy integration |
| [imserver-sdk-java](https://github.com/juggleim/imserver-sdk-java) | Java version of imserver-sdk |

> The desktop imsdk-pc is not yet open-sourced — contact support for details.

## 🚀 Quick Deploy

For the fastest experience, follow the official one-click deployment guide:

👉 **[One-Click Deployment Guide](https://www.juggle.im/docs/guide/deploy/quickdeploy/)**

## 🛠 Manual Deployment

<details>
<summary>Click to expand the full manual deployment steps</summary>

### 1. Install and Initialize MySQL

Create the database schema:
```sql
CREATE SCHEMA `jim_db`;
```

Initialize the table structure (the SQL file lives at `sql/imserver.sql`):
```bash
mysql -u{db_user} -p{db_password} jim_db < sql/imserver.sql
```

### 2. Install MongoDB (optional)

Only required when using MongoDB to store message data (`msgStoreEngine: mongo`).

### 3. Start im-server

The working directory is `im-server/launcher`, where `conf` holds config files and `logs` is the runtime log directory.

**Edit the config file** `im-server/launcher/conf/config.yml`:
```yaml
defaultPort: 9003       # im-server default listening port
nodeName: testNode      # node name
nodeHost: 127.0.0.1     # node IP
msgStoreEngine: mysql   # message store engine: mysql (default) or mongo

log:
  logPath: ./logs       # runtime log directory
  logName: jim-info     # runtime log filename prefix
  visual: false         # enable visual logs (write to a KV store, queryable in the admin console)

mysql:                  # MySQL configuration
  user: root
  password: 123456
  address: 127.0.0.1:3306
  name: im_db

# mongodb:              # MongoDB config, active when msgStoreEngine is "mongo"
#   address: 127.0.0.1:27017
#   name: jim_msgs

# apiGateway:           # server-side API port for business app servers; defaults to defaultPort
#   httpPort: 9001

# connectManager:       # long-connection port; defaults to defaultPort
#   wsPort: 9003

adminGateway:           # built-in admin console, default credentials admin/123456
  httpPort: 8090
```

**Start the service** from the `im-server/launcher` directory:
```bash
go run main.go
```

### 4. Configure Public Access Addresses

Ports that need to be exposed:

| Port | Protocol | Description |
| ---: | :---: | :--- |
| 9003 | http | Server-side API port, called by business servers (e.g. jugglechat-server) |
| 9003 | websocket | IM long-connection port for client SDKs |
| 8090 | http | Admin console port, default credentials admin/123456 |

Configure exposure however suits your environment (public IP, Nginx reverse proxy, load balancer, etc.). For local testing, an intranet IP is enough.

**Register the long-connection address** by inserting a config row into the database:
```sql
insert into globalconfs (conf_key, conf_value) values ('connect_address', '{"default":["127.0.0.1:9002"]}');
```
Replace `127.0.0.1` with your machine's intranet IP or public IP/domain. This address is delivered to client SDKs by the navigator service.

</details>

## 🏢 Create an Application (Tenant)

JuggleIM is a **multi-tenant** system — a single deployment can host multiple appkeys (tenants) with fully isolated data.

**Create a tenant via the admin API** (`app_key` is the tenant identifier and must be unique):
```bash
curl --request POST \
  --url http://127.0.0.1:8090/admingateway/apps/create \
  --data '{
    "app_key":"appkey",
    "app_name":"appname"
}'
```

Example response:
```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "app_name": "appname",
        "app_key": "appkey",
        "app_secret": "hciKcc6sXRDjYUQp"
    }
}
```

You can also log in to the admin console at `http://127.0.0.1:8090` (default credentials `admin/123456`) to view and manage your applications.

## 🔌 Business Server / Client Integration

**1) Business Server Integration**

| Item | Example | Notes |
| ---: | :---: | :--- |
| IM server-side API address | `http://127.0.0.1:9003` | Used by your business server to call IM APIs (register users, create groups, send system messages, etc.). See the [API Reference](https://www.juggle.im/docs/server/api/) |
| app_key | `appkey1` | Tenant identifier, unique within the system |
| app_secret | `hciKcc6sXRDjYUQp` | Auth secret generated on app creation (must be 16 chars if custom). **Use only on the business server — never expose it to clients** |

**2) Client SDK Integration**

| Item | Example | Notes |
| ---: | :---: | :--- |
| IM connection address | `ws://127.0.0.1:9003` | Passed to the client SDK on init. See [Quick Start](https://www.juggle.im/docs/client/quickstart/android/) |
| app_key | `appkey1` | Must match the value used on the business server |

## 💬 Community

Interested in IM or have integration questions? Join the community and let's chat 👇

- [Telegram Group (Chinese)](https://t.me/juggleim_zh)
- [Open an Issue](https://github.com/juggleim/im-server/issues)

## 🤝 Contributing

Contributions of all kinds are welcome! You can:

- Open an [Issue](https://github.com/juggleim/im-server/issues) to report bugs or request features
- Submit a Pull Request to improve the code or docs
- Share the projects you build on top of JuggleIM

## ⭐ Star History

If JuggleIM has helped you, please give us a Star — your support drives our continued development!

[![Star History Chart](https://api.star-history.com/svg?repos=juggleim/im-server&type=Date)](https://star-history.com/#juggleim/im-server&Date)

## 📄 License

This project is released under the [LICENSE](./LICENSE).
