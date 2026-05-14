#!/usr/bin/env bash
set -euo pipefail

API_DOMAIN="clothapi.235235.vip"
WEB_DOMAIN="clothweb.235235.vip"
H5_DOMAIN="235235.vip"
API_LOCAL_PORT="8899"

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

days_to_expire() {
  local host="$1"
  local end
  end=$(echo | openssl s_client -servername "$host" -connect "$host:443" 2>/dev/null | openssl x509 -noout -enddate | cut -d= -f2)
  [[ -n "$end" ]] || { echo "-1"; return; }

  local end_ts now_ts
  end_ts=$(date -d "$end" +%s)
  now_ts=$(date +%s)
  echo $(( (end_ts - now_ts) / 86400 ))
}

echo "==== RichRoad BT Deployment Check ===="

need_cmd curl
need_cmd openssl
need_cmd ss

if command -v systemctl >/dev/null 2>&1; then
  if systemctl is-active --quiet richroad; then
    ok "Service richroad is running"
  else
    warn "Service richroad is not active (or service name differs)"
  fi
fi

if ss -lntp | grep -q ":${API_LOCAL_PORT} "; then
  ok "Backend is listening on :${API_LOCAL_PORT}"
else
  fail "Backend is NOT listening on :${API_LOCAL_PORT}"
fi

api_health=$(http_code "https://${API_DOMAIN}/health")
[[ "$api_health" == "200" ]] && ok "API health: 200" || fail "API health check failed, got ${api_health}"

api_swagger=$(http_code "https://${API_DOMAIN}/swagger/index.html")
[[ "$api_swagger" == "200" || "$api_swagger" == "301" || "$api_swagger" == "302" ]] \
  && ok "API swagger reachable (${api_swagger})" \
  || warn "API swagger status: ${api_swagger}"

web_home=$(http_code "https://${WEB_DOMAIN}")
[[ "$web_home" == "200" || "$web_home" == "301" || "$web_home" == "302" ]] \
  && ok "Web admin reachable (${web_home})" \
  || fail "Web admin unreachable, status ${web_home}"

h5_home=$(http_code "https://${H5_DOMAIN}")
[[ "$h5_home" == "200" || "$h5_home" == "301" || "$h5_home" == "302" ]] \
  && ok "Uni H5 reachable (${h5_home})" \
  || fail "Uni H5 unreachable, status ${h5_home}"

cors_h5=$(curl -k -sSI -X OPTIONS \
  -H "Origin: https://${H5_DOMAIN}" \
  -H "Access-Control-Request-Method: POST" \
  "https://${API_DOMAIN}/health" | tr -d '\r' | awk -F': ' 'tolower($1)=="access-control-allow-origin"{print $2}' | tail -n1)

if [[ "$cors_h5" == "https://${H5_DOMAIN}" ]]; then
  ok "CORS allows H5 origin"
else
  warn "CORS for H5 origin not returned as expected (got: ${cors_h5:-<empty>})"
fi

cors_web=$(curl -k -sSI -X OPTIONS \
  -H "Origin: https://${WEB_DOMAIN}" \
  -H "Access-Control-Request-Method: POST" \
  "https://${API_DOMAIN}/health" | tr -d '\r' | awk -F': ' 'tolower($1)=="access-control-allow-origin"{print $2}' | tail -n1)

if [[ "$cors_web" == "https://${WEB_DOMAIN}" ]]; then
  ok "CORS allows Web origin"
else
  warn "CORS for Web origin not returned as expected (got: ${cors_web:-<empty>})"
fi

for host in "$API_DOMAIN" "$WEB_DOMAIN" "$H5_DOMAIN"; do
  d=$(days_to_expire "$host")
  if [[ "$d" == "-1" ]]; then
    warn "Unable to parse TLS cert for ${host}"
  elif (( d < 7 )); then
    warn "TLS cert for ${host} expires in ${d} day(s)"
  else
    ok "TLS cert for ${host} valid for ${d} day(s)"
  fi
done

echo "==== Check Complete ===="
