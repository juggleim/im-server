#!/usr/bin/env bash
set -euo pipefail

for command_name in curl jq openssl; do
  if ! command -v "${command_name}" >/dev/null 2>&1; then
    echo "Required command not found: ${command_name}" >&2
    exit 1
  fi
done

admin_base_url="${ADMIN_BASE_URL:-http://127.0.0.1:8090/admingateway}"
api_base_url="${API_BASE_URL:-http://127.0.0.1:9001/apigateway}"
admin_account="${ADMIN_ACCOUNT:-admin}"
admin_password="${ADMIN_PASSWORD:-123456}"
run_suffix="$(date +%s)"
app_key="${APP_KEY:-quickstart${run_suffix}}"
alice_id="alice_${run_suffix}"
bob_id="bob_${run_suffix}"

require_success() {
  local operation="$1"
  local response="$2"
  local code
  code="$(jq -r '.code // empty' <<<"${response}")"
  if [[ "${code}" != "0" ]]; then
    echo "${operation} failed: ${response}" >&2
    exit 1
  fi
}

login_response="$(curl --fail --silent --show-error \
  --request POST "${admin_base_url}/login" \
  --header 'Content-Type: application/json' \
  --data "$(jq -nc \
    --arg account "${admin_account}" \
    --arg password "${admin_password}" \
    '{account: $account, password: $password}')")"
require_success "Admin login" "${login_response}"
admin_token="$(jq -er '.data.authorization' <<<"${login_response}")"

app_response="$(curl --fail --silent --show-error \
  --request POST "${admin_base_url}/apps/create" \
  --header 'Content-Type: application/json' \
  --header "Authorization: ${admin_token}" \
  --data "$(jq -nc \
    --arg app_key "${app_key}" \
    '{app_key: $app_key, app_name: "Server API Quick Start"}')")"
require_success "Create application" "${app_response}"
app_secret="$(jq -er '.data.app_secret' <<<"${app_response}")"

signed_post() {
  local endpoint="$1"
  local body="$2"
  local nonce timestamp signature
  nonce="$(openssl rand -hex 12)"
  timestamp="$(date +%s)000"
  signature="$(printf '%s' "${app_secret}${nonce}${timestamp}" | openssl dgst -sha1 | awk '{print $NF}')"

  curl --fail --silent --show-error \
    --request POST "${api_base_url}${endpoint}" \
    --header 'Content-Type: application/json' \
    --header "appkey: ${app_key}" \
    --header "nonce: ${nonce}" \
    --header "timestamp: ${timestamp}" \
    --header "signature: ${signature}" \
    --data "${body}"
}

alice_response="$(signed_post '/users/register' "$(jq -nc \
  --arg user_id "${alice_id}" \
  '{user_id: $user_id, nickname: "Alice"}')")"
require_success "Register Alice" "${alice_response}"
alice_token="$(jq -er '.data.token' <<<"${alice_response}")"

bob_response="$(signed_post '/users/register' "$(jq -nc \
  --arg user_id "${bob_id}" \
  '{user_id: $user_id, nickname: "Bob"}')")"
require_success "Register Bob" "${bob_response}"
bob_token="$(jq -er '.data.token' <<<"${bob_response}")"

message_response="$(signed_post '/messages/private/send' "$(jq -nc \
  --arg sender_id "${alice_id}" \
  --arg receiver_id "${bob_id}" \
  --arg content '{"content":"Hello from the server API quick start"}' \
  '{sender_id: $sender_id, receiver_id: $receiver_id, msg_type: "jg:text", msg_content: $content}')")"
require_success "Send private message" "${message_response}"
message_id="$(jq -er '.data[0].msg_id' <<<"${message_response}")"

cat <<OUTPUT
Server API quick start completed.
Application: ${app_key}
Alice user ID: ${alice_id}
Alice token: ${alice_token}
Bob user ID: ${bob_id}
Bob token: ${bob_token}
Message ID: ${message_id}

Keep the application secret and server API signature logic on your trusted server.
The user tokens above are for this local test environment only.
OUTPUT
