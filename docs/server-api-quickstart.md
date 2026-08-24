# Server API quick start

This guide takes a fresh local Docker Compose deployment through the minimum server-side flow needed before a client SDK can connect:

1. Sign in to the local admin API.
2. Create an application (tenant).
3. Register two synthetic users and obtain their client tokens.
4. Send a private message through the server API.

The example was verified against the current [`docker-compose.yml`](../docker-compose.yml). It uses only generated local identifiers and message content.

## Prerequisites

Install Docker Compose, `curl`, `jq`, and `openssl`, then start the local stack:

```bash
docker compose up -d --wait --wait-timeout 120
docker compose ps
```

MySQL should be `healthy`, and `im-server` should be running. If not, use the [Docker Compose troubleshooting guide](./docker-troubleshooting.md).

## Run the verified example

```bash
bash examples/server-api-quickstart.sh
```

The script uses the local-only admin credentials from `docker-compose.yml`, creates a unique application, registers Alice and Bob, and sends a stored `jg:text` private message from Alice to Bob. It prints the generated user IDs, user tokens, and message ID, but never prints the application secret.

You can override the local endpoints and admin credentials without editing the script:

```bash
ADMIN_BASE_URL=http://127.0.0.1:8090/admingateway \
API_BASE_URL=http://127.0.0.1:9001/apigateway \
ADMIN_ACCOUNT=admin \
ADMIN_PASSWORD='<local-admin-password>' \
bash examples/server-api-quickstart.sh
```

## APIs used

| Step | Endpoint | Result |
| :--- | :--- | :--- |
| Admin login | `POST /admingateway/login` | Short-lived local admin authorization |
| Create application | `POST /admingateway/apps/create` | Application key and server-only application secret |
| [Register user](https://www.juggle.im/docs/server/user/register/) | `POST /apigateway/users/register` | User ID and client token |
| [Send private message](https://www.juggle.im/docs/server/message/privatemsg/) | `POST /apigateway/messages/private/send` | Message ID for each receiver |

Server API calls use the documented [signature headers](https://www.juggle.im/docs/server/api/#header): `appkey`, `nonce`, `timestamp`, and `signature`. The signature is the SHA-1 digest of:

```text
app_secret + nonce + timestamp
```

The example recomputes the signature for every request. In a real integration, your business server performs this work and then returns only the user token and connection information required by the client SDK.

## Connect a client SDK

Use all three values together:

- the generated application key;
- one generated user ID and its matching token;
- the local WebSocket address `ws://127.0.0.1:9003`.

Follow the relevant [client SDK quick start](https://www.juggle.im/docs/client/quickstart/android/) for platform-specific initialization. Do not mix a token from one application or user with another application or user ID.

## Security boundary

Keep these values and operations on a trusted server:

- admin credentials and admin authorization;
- application secret;
- server API signing logic;
- calls that create applications, register users, or send privileged server messages.

A client application may receive its own application key, user ID, short-lived user token, and navigator/WebSocket address. Never ship the application secret or admin authorization in web, desktop, or mobile client code. Never log them or commit them to source control.

The script prints synthetic local user tokens so they can be copied into a local SDK test. Treat production user tokens as credentials and deliver them only to the authenticated user they belong to.
