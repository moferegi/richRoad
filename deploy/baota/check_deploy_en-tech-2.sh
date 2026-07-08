#!/usr/bin/env bash
set -euo pipefail

API_DOMAIN="clothapi-en-tech-2.example.com"
API_LOCAL_PORT="8658"
SERVICE_NAME="richroad-en-tech-2"

ok()   { printf "[OK] %s\n" "$1"; }
warn() { printf "[WARN] %s\n" "$1"; }
fail() { printf "[FAIL] %s\n" "$1"; exit 1; }

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "Missing command: $1"
}

http_code() {
  local url="$1"
  curl -k -sS -o /dev/null -w "%{http_code}" "$url"
}

echo "==== RichRoad en-tech-2 Deployment Check ===="

need_cmd curl
need_cmd ss

if command -v systemctl >/dev/null 2>&1; then
  if systemctl is-active --quiet "${SERVICE_NAME}"; then
    ok "Service ${SERVICE_NAME} is running"
  else
    warn "Service ${SERVICE_NAME} is not active"
  fi
fi

if ss -lntp | grep -q ":${API_LOCAL_PORT} "; then
  ok "Backend is listening on :${API_LOCAL_PORT}"
else
  fail "Backend is NOT listening on :${API_LOCAL_PORT}"
fi

api_health=$(http_code "https://${API_DOMAIN}/health")
[[ "$api_health" == "200" ]] && ok "API health: 200" || warn "API health check status: ${api_health}"

echo "==== Check Complete ===="
