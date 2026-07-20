#!/usr/bin/env bash
set -euo pipefail

for command_name in curl docker go jq openssl shasum; do
  if ! command -v "${command_name}" >/dev/null 2>&1; then
    echo "Required command not found: ${command_name}" >&2
    exit 1
  fi
done

scenario="${1:-all}"
if [[ "${scenario}" != "all" && "${scenario}" != "private" && "${scenario}" != "group" ]]; then
  echo "Usage: $0 [all|private|group]" >&2
  exit 1
fi

clients="${BENCH_CLIENTS:-50}"
rate="${BENCH_RATE:-200}"
warmup="${BENCH_WARMUP:-10s}"
duration="${BENCH_DURATION:-30s}"
delivery_grace="${BENCH_DELIVERY_GRACE:-3s}"
payload_bytes="${BENCH_PAYLOAD_BYTES:-256}"
group_senders="${BENCH_GROUP_SENDERS:-2}"
output_dir="${BENCH_OUTPUT_DIR:-benchmark-results}"
export MYSQL_HOST_PORT="${MYSQL_HOST_PORT:-13306}"
export API_HOST_PORT="${API_HOST_PORT:-19001}"
export NAV_HOST_PORT="${NAV_HOST_PORT:-19002}"
export WS_HOST_PORT="${WS_HOST_PORT:-19003}"
export ADMIN_HOST_PORT="${ADMIN_HOST_PORT:-18090}"
export PPROF_HOST_PORT="${PPROF_HOST_PORT:-16060}"
admin_base_url="${ADMIN_BASE_URL:-http://127.0.0.1:${ADMIN_HOST_PORT}/admingateway}"
api_base_url="${API_BASE_URL:-http://127.0.0.1:${API_HOST_PORT}/apigateway}"
ws_url="${WS_URL:-ws://127.0.0.1:${WS_HOST_PORT}}"
admin_account="${ADMIN_ACCOUNT:-admin}"
admin_password="${ADMIN_PASSWORD:-123456}"
run_timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
app_key="bench$(date +%s)"

echo "Building and starting the local Docker Compose stack..."
docker compose up --detach --build

login_response=""
for _ in $(seq 1 60); do
  if login_response="$(curl --fail --silent --show-error \
    --request POST "${admin_base_url}/login" \
    --header 'Content-Type: application/json' \
    --data "$(jq -nc --arg account "${admin_account}" --arg password "${admin_password}" \
      '{account: $account, password: $password}')" 2>/dev/null)"; then
    if [[ "$(jq -r '.code // -1' <<<"${login_response}")" == "0" ]]; then
      break
    fi
  fi
  sleep 2
done

if [[ -z "${login_response}" || "$(jq -r '.code // -1' <<<"${login_response}")" != "0" ]]; then
  echo "Local stack did not become ready: ${login_response:-no response}" >&2
  exit 1
fi

admin_token="$(jq -er '.data.authorization' <<<"${login_response}")"
app_response="$(curl --fail --silent --show-error \
  --request POST "${admin_base_url}/apps/create" \
  --header 'Content-Type: application/json' \
  --header "Authorization: ${admin_token}" \
  --data "$(jq -nc --arg app_key "${app_key}" \
    '{app_key: $app_key, app_name: "Reproducible Benchmark"}')")"
if [[ "$(jq -r '.code // -1' <<<"${app_response}")" != "0" ]]; then
  echo "Create benchmark application failed: ${app_response}" >&2
  exit 1
fi

export JIM_BENCH_APP_KEY="${app_key}"
export JIM_BENCH_APP_SECRET
JIM_BENCH_APP_SECRET="$(jq -er '.data.app_secret' <<<"${app_response}")"
sleep 2
export JIM_BENCH_SERVER_COMMIT="$(git rev-parse HEAD)"
compose_sha="$(shasum -a 256 docker-compose.yml | awk '{print $1}')"
export JIM_BENCH_ENVIRONMENT="Local Docker Compose; server and load generator share the host; no explicit container CPU or memory limits; compose_sha256=${compose_sha}"
export JIM_BENCH_DATABASE="MySQL 8.0 Docker image; utf8mb4_0900_ai_ci; default Compose configuration"
mkdir -p "${output_dir}"

run_scenario() {
  local selected_scenario="$1"
  echo "Running ${selected_scenario} workload..."
  go run ./cmd/jimbench \
    --scenario "${selected_scenario}" \
    --ws-url "${ws_url}" \
    --api-url "${api_base_url}" \
    --clients "${clients}" \
    --group-senders "${group_senders}" \
    --rate "${rate}" \
    --warmup "${warmup}" \
    --duration "${duration}" \
    --delivery-grace "${delivery_grace}" \
    --payload-bytes "${payload_bytes}" \
    --output "${output_dir}/${run_timestamp}-${selected_scenario}.json"
}

if [[ "${scenario}" == "all" || "${scenario}" == "private" ]]; then
  run_scenario private
fi
if [[ "${scenario}" == "all" || "${scenario}" == "group" ]]; then
  run_scenario group
fi

echo "Benchmark results are in ${output_dir}/."
echo "The local stack remains running for inspection. Stop it with: docker compose down"
