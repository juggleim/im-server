# APNs P8 runtime support

This standalone server adds APNs P8 provider-key support without moving legacy P12 certificates into the P8 manager.

## Database contract

Migration `commons/dbcommons/sqls/20260908.sql` adds these fields to `ioscertificates`:

- `auth_type`: `p12` or `p8`; existing rows default to `p12`.
- `p8_key_id`: Apple Key ID.
- `p8_team_id`: Apple Team ID.
- `p8_private_key`: original PKCS#8 PEM file bytes, stored directly without application-layer encryption.
- `p8_key_name`: display-only upload filename retained as part of the console/database contract; the runtime does not use it.
- `config_version`: configuration refresh version.

The baseline schema in `sql/imserver.sql` contains the same fields. Database users and backup operators can read the raw P8 key and must protect access accordingly.

## Runtime behavior

- P8 requires an unencrypted PKCS#8 PEM ECDSA P-256 key, valid Key ID and Team ID.
- P8 ordinary and VoIP requests share one APNs token client; topics remain the bundle ID and `<bundle>.voip`.
- The P8 manager accepts only rows whose `auth_type` is exactly `p8`; its config contains only identity, environment, P8 credentials and version fields.
- P12 remains on the original `GetIosPushConf`/`initIosPushConf` path, including its five-minute cache and `Client.Push` 60-second behavior.
- The original P12 initialization semantics are unchanged: a successfully parsed ordinary certificate returns immediately, so when both certificates exist the ordinary client is used as the actual fallback; VoIP initialization is reached only when ordinary initialization does not return.
- P8 configuration is checked every five minutes on the next P8 send. Unchanged configurations retain the APNs connection and provider JWT.
- `IM_APNS_PUSH_TIMEOUT_SECONDS` optionally sets the total operation timeout to 1–120 seconds; the default is 10 seconds.
- The runtime does not read an APNs encryption master key and does not contain control-plane upload APIs.

## Administration dependency

This repository intentionally does not copy the management API or frontend. They are provided by the `github.com/juggleim/imserver-console` dependency. The deployed console module version must support the six database fields and raw `p8_file` uploads. The currently pinned module version should be updated after that console implementation is published; do not add a workstation-local `replace` directive for production delivery.

## Targeted verification

```bash
go test ./services/pushmanager/services/apnspush ./services/pushmanager/services ./services/pushmanager -count=1
go test -race ./services/pushmanager/services/apnspush -count=1
go test -race ./services/pushmanager ./services/pushmanager/services -run '^TestAPNS' -count=1
go build ./...
```

Some repository integration tests require external databases or push vendors and are not part of the offline P8 verification.
