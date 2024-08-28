#!/usr/bin/env bash

set -e
set -o noglob

info() {
    echo '[INFO] ' "$@"
}

warn() {
    echo '[WARN] ' "$@" >&2
}

fatal() {
    echo '[ERROR] ' "$@" >&2
    exit 1
}

#BASE_DIR="/opt/conf"
CONF_DIR="/opt/conf"

BASE_DIR="/opt"

CONFIG_FILE="config.yml"
CONFIG_FILE_TEMPLATE="config_template.yaml"


env_vars=(
  CLUSTER_NAME
  POD_NAME
  POD_IP
  RPC_PORT
  MSG_STORE_ENGINE
  TAG_STORE_ENGINE
  MYSQL_ROOT_PASSWORD
  MYSQL_ADDR
  MYSQL_DB_NAME
  MONGODB_ADDR
  MONGODB_ROOT_PASSWORD
  WS_PORT
  PROXY_PORT
  API_HTTP_PORT
  NAV_HTTP_PORT
  ADMIN_HTTP_PORT
)

# 检查环境变量是否存在
function check_env_var() {
  local var_name="$1"

  if [ -z "${!var_name}" ]; then
    warn "$var_name environment variable is not set. Skipping replacement."
    return 1
  fi
}

# 替换文件中的变量
function replace_env_var() {
  local var_name="$1"
  local config_file="$2"
  local var_value="${!var_name}"

  if ! check_env_var "$var_name"; then
    return
  fi

  sed -i "s/{{ $var_name }}/$var_value/g" "$config_file"
}

function init_config() {
    for var_name in "${env_vars[@]}"; do
      check_env_var "$var_name"
      replace_env_var "$var_name" "$BASE_DIR/$CONFIG_FILE_TEMPLATE"
    done
    cp "$BASE_DIR/$CONFIG_FILE_TEMPLATE" "$CONF_DIR/$CONFIG_FILE"
}

init_config

/opt/imserver 2>&1 | tee logs/jim_console.log