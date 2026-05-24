# Security Hardening Runbook

This runbook documents rollout, verification, and rollback steps for the recent security hardening changes.

## Scope

The runbook covers:
- Casbin legacy over-grant cleanup switch: `CS_CLEANUP_LEGACY_CASBIN_OVERGRANT`
- Export download token transport switch: `CS_EXPORT_ALLOW_QUERY_TOKEN`
- Upload size guardrail: `CS_UPLOAD_MAX_SIZE_MB`
- Scan upload ticket TTL switch: `CS_SCAN_UPLOAD_TICKET_TTL_SECONDS`
- AutoCode SQL identifier hardening (MSSQL/SQLite/DropTable)
- Web/Uni WebSocket query-token fallback switch: `VITE_WS_ALLOW_QUERY_TOKEN`

Sensitive API ignore cleanup:
- `/init/initdb` and `/init/checkdb` are no longer recommended in `sys_ignore_apis`.
- `/system/reloadSystem` is also no longer recommended in `sys_ignore_apis`.
- Startup now attempts to remove these sensitive ignore entries automatically to reduce audit blind spots.

API-layer role allowlist defense-in-depth:
- `POST /system/reloadSystem` and `POST /system/getServerInfo` are guarded at API layer (role `888`) in addition to Casbin.
- `POST /system/getSystemConfig` and `POST /system/setSystemConfig` are guarded at API layer (roles `888`, `8881`, `9528`) in addition to Casbin.
- File library sensitive APIs are guarded at API layer (roles `888`, `8881`, `9528`) in addition to Casbin:
  - `POST /fileUploadAndDownload/deleteFile`
  - `POST /fileUploadAndDownload/editFileName`
  - `POST /fileUploadAndDownload/getFileList`
  - `POST /fileUploadAndDownload/importURL`
  - `GET /fileUploadAndDownload/listFolders`
- `POST /fileUploadAndDownload/signURL` now normalizes and validates `filePath` (rejects empty, traversal-like, query/fragment paths) before signing.

Hotlink permissions self-heal:
- Startup registration now treats `GET /fileUploadAndDownload/listFolders` as admin-only (`888`, `8881`, `9528`).
- Legacy over-grant for role `8080` on `GET /fileUploadAndDownload/listFolders` is removed automatically on startup.

Upload delete safety guard:
- `DeleteFile` now validates and normalizes object-storage keys before deletion.
- If a record has empty or illegal `key`, backend skips object deletion and only removes DB record to prevent risky delete calls against storage prefixes.

Upload input boundary hardening:
- `EditFileName` now rejects invalid file IDs and empty/oversized names.
- `GetFileList` enforces safe pagination defaults and maximum page size to reduce query amplification.
- `ImportURL` now enforces batch-size limits and URL validation (`http/https` or absolute path), and sanitizes imported name/keywords fields.

File ownership isolation:
- `exa_file_upload_and_downloads` now stores `created_by` (creator user ID).
- Upload and import operations write `created_by` from the current operator.
- File management operations (`deleteFile`, `editFileName`, `getFileList`) are scoped by ownership for non-super roles.
- Role `888` keeps global management scope; non-super file managers can only operate on records created by themselves.
- Legacy rows without owner information (`created_by=0`) should be handled by role `888` or backfilled before delegating to non-super managers.

Attack observability:
- When `POST /sysError/createSysError` hits rate limit, backend records attack type `sys_error_rate_limit` into `sys_attack_logs` and surfaces it in `/sysBannedIP/getAttackStats` as `sysErrorRate`.
- Auto-ban counters are tracked per attack type; `sys_error_rate_limit` has an independent threshold (`security_attack_auto_ban_threshold_sys_error_rate_limit`) and no longer mixes with login-fail counters.
- `POST /fileUploadAndDownload/uploadByTicket` now records:
  - `scan_upload_ticket_rate_limit` (shown as `scanTicketRate` in attack stats)
  - `scan_upload_ticket_invalid` (shown as `scanTicketInvalid` in attack stats)
- Scan-upload ticket attack types have independent auto-ban thresholds:
  - `security_attack_auto_ban_threshold_scan_upload_ticket_rate_limit`
  - `security_attack_auto_ban_threshold_scan_upload_ticket_invalid`

AutoCode SQL hardening:
- MSSQL metadata queries now normalize database identifiers and use bound parameters for table names.
- SQLite metadata query (`PRAGMA table_info`) now quotes identifiers safely.
- AutoCode rollback `DropTable` now enforces strict table-name validation before executing DDL.

WebSocket token transport hardening:
- Web/Uni client fallback from subprotocol to query token is disabled by default.
- Only when `VITE_WS_ALLOW_QUERY_TOKEN=true` is explicitly configured should query-token fallback be used as temporary compatibility.

Scan upload token transport hardening:
- Web admin QR upload now uses one-time short-lived ticket and no longer places long-lived JWT in URL query.
- New ticket flow:
  - `POST /fileUploadAndDownload/createScanUploadTicket` (issue ticket, requires JWT-authenticated file manager role)
  - `POST /fileUploadAndDownload/uploadByTicket` (consume one-time ticket on upload, ticket only accepted from form-data body)
- Ticket lifetime defaults to 180 seconds and can be configured by `security_scan_upload_ticket_ttl_seconds` / `CS_SCAN_UPLOAD_TICKET_TTL_SECONDS`.
- Ticket consume is one-time by design (Redis `GET+DEL`; local in-memory fallback when Redis unavailable).
- `uploadByTicket` endpoint now applies upload rate-limit by `IP + User-Agent` window (Redis first, local-memory fallback).

## Runtime Config Priority

Security runtime settings now support backend management.

Priority order:
1. `client_sys_config.config_value` (admin panel)
2. Environment variable fallback
3. Built-in default value

### Backend Keys (group: `security`)

- `security_upload_max_size_mb`
- `security_upload_strict_validation_enabled`
- `security_scan_upload_ticket_ttl_seconds`
- `security_scan_upload_ticket_upload_rate_limit_per_ip`
- `security_scan_upload_ticket_upload_rate_limit_window_seconds`
- `security_export_allow_query_token`
- `security_ws_allow_query_token`
- `security_ws_max_conns_per_ip`
- `captcha_rate_limit`
- `captcha_rate_limit_window_seconds`
- `security_login_ip_rate_limit_per_minute`
- `security_login_ip_rate_limit_window_seconds`
- `security_init_api_enabled`
- `security_init_api_private_network_only`
- `security_trusted_proxies`
- `security_public_rate_limit_enabled`
- `security_public_rate_limit_window_seconds`
- `security_public_rate_limit_max_requests`
- `security_attack_auto_ban_enabled`
- `security_attack_auto_ban_window_seconds`
- `security_attack_auto_ban_duration_minutes`
- `security_attack_auto_ban_threshold_default`
- `security_attack_auto_ban_threshold_sys_error_rate_limit`
- `security_attack_auto_ban_threshold_scan_upload_ticket_rate_limit`
- `security_attack_auto_ban_threshold_scan_upload_ticket_invalid`
- `security_sys_error_create_rate_limit_per_minute`
- `security_sys_error_create_rate_limit_window_seconds`
- `security_visitor_heartbeat_rate_limit_per_minute`
- `security_visitor_heartbeat_dedupe_seconds`
- `security_tryon_create_rate_limit_per_minute`
- `security_tryon_create_concurrency_limit`
- `security_cleanup_legacy_casbin_overgrant`

## Recommended Profiles

Apply one of the following profiles based on your environment.

### Profile A: Production (security-first)

- `security_upload_max_size_mb = 20`
- `security_upload_strict_validation_enabled = true`
- `security_scan_upload_ticket_ttl_seconds = 180`
- `security_scan_upload_ticket_upload_rate_limit_per_ip = 30`
- `security_scan_upload_ticket_upload_rate_limit_window_seconds = 60`
- `security_export_allow_query_token = false`
- `security_ws_allow_query_token = false`
- `security_ws_max_conns_per_ip = 8`
- `captcha_rate_limit = 10`
- `captcha_rate_limit_window_seconds = 60`
- `security_login_ip_rate_limit_per_minute = 30`
- `security_login_ip_rate_limit_window_seconds = 60`
- `security_init_api_enabled = false`
- `security_init_api_private_network_only = true`
- `security_trusted_proxies = 127.0.0.1,::1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7`
- `security_public_rate_limit_enabled = true`
- `security_public_rate_limit_window_seconds = 60`
- `security_public_rate_limit_max_requests = 300`
- `security_attack_auto_ban_enabled = true`
- `security_attack_auto_ban_window_seconds = 3600`
- `security_attack_auto_ban_duration_minutes = 60`
- `security_attack_auto_ban_threshold_default = 10`
- `security_attack_auto_ban_threshold_sys_error_rate_limit = 30`
- `security_attack_auto_ban_threshold_scan_upload_ticket_rate_limit = 30`
- `security_attack_auto_ban_threshold_scan_upload_ticket_invalid = 50`
- `security_sys_error_create_rate_limit_per_minute = 30`
- `security_sys_error_create_rate_limit_window_seconds = 60`
- `security_visitor_heartbeat_rate_limit_per_minute = 120`
- `security_visitor_heartbeat_dedupe_seconds = 3`
- `security_tryon_create_rate_limit_per_minute = 20`
- `security_tryon_create_concurrency_limit = 2`
- `security_cleanup_legacy_casbin_overgrant = false`

Recommended for public traffic and long-term stable operation.

### Profile B: Integration/Legacy Compatibility (temporary)

- `security_upload_max_size_mb = 20`
- `security_upload_strict_validation_enabled = true`
- `security_scan_upload_ticket_ttl_seconds = 300`
- `security_scan_upload_ticket_upload_rate_limit_per_ip = 60`
- `security_scan_upload_ticket_upload_rate_limit_window_seconds = 60`
- `security_export_allow_query_token = true`
- `security_ws_allow_query_token = true`
- `security_ws_max_conns_per_ip = 16`
- `captcha_rate_limit = 20`
- `captcha_rate_limit_window_seconds = 60`
- `security_login_ip_rate_limit_per_minute = 60`
- `security_login_ip_rate_limit_window_seconds = 60`
- `security_init_api_enabled = true`
- `security_init_api_private_network_only = true`
- `security_trusted_proxies = 127.0.0.1,::1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7`
- `security_public_rate_limit_enabled = true`
- `security_public_rate_limit_window_seconds = 60`
- `security_public_rate_limit_max_requests = 600`
- `security_attack_auto_ban_enabled = true`
- `security_attack_auto_ban_window_seconds = 3600`
- `security_attack_auto_ban_duration_minutes = 60`
- `security_attack_auto_ban_threshold_default = 20`
- `security_attack_auto_ban_threshold_sys_error_rate_limit = 60`
- `security_attack_auto_ban_threshold_scan_upload_ticket_rate_limit = 60`
- `security_attack_auto_ban_threshold_scan_upload_ticket_invalid = 100`
- `security_sys_error_create_rate_limit_per_minute = 120`
- `security_sys_error_create_rate_limit_window_seconds = 60`
- `security_visitor_heartbeat_rate_limit_per_minute = 240`
- `security_visitor_heartbeat_dedupe_seconds = 2`
- `security_tryon_create_rate_limit_per_minute = 40`
- `security_tryon_create_concurrency_limit = 4`
- `security_cleanup_legacy_casbin_overgrant = false`

Use this only when old clients are not yet upgraded. Switch back to Profile A after upgrade.

### SQL Example (batch update)

```sql
UPDATE client_sys_config
SET config_value = CASE config_key
  WHEN 'security_upload_max_size_mb' THEN '20'
  WHEN 'security_upload_strict_validation_enabled' THEN 'true'
  WHEN 'security_scan_upload_ticket_ttl_seconds' THEN '180'
  WHEN 'security_scan_upload_ticket_upload_rate_limit_per_ip' THEN '30'
  WHEN 'security_scan_upload_ticket_upload_rate_limit_window_seconds' THEN '60'
  WHEN 'security_export_allow_query_token' THEN 'false'
  WHEN 'security_ws_allow_query_token' THEN 'false'
  WHEN 'security_ws_max_conns_per_ip' THEN '8'
  WHEN 'captcha_rate_limit' THEN '10'
  WHEN 'captcha_rate_limit_window_seconds' THEN '60'
  WHEN 'security_login_ip_rate_limit_per_minute' THEN '30'
  WHEN 'security_login_ip_rate_limit_window_seconds' THEN '60'
  WHEN 'security_init_api_enabled' THEN 'false'
  WHEN 'security_init_api_private_network_only' THEN 'true'
  WHEN 'security_trusted_proxies' THEN '127.0.0.1,::1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7'
  WHEN 'security_public_rate_limit_enabled' THEN 'true'
  WHEN 'security_public_rate_limit_window_seconds' THEN '60'
  WHEN 'security_public_rate_limit_max_requests' THEN '300'
  WHEN 'security_attack_auto_ban_enabled' THEN 'true'
  WHEN 'security_attack_auto_ban_window_seconds' THEN '3600'
  WHEN 'security_attack_auto_ban_duration_minutes' THEN '60'
  WHEN 'security_attack_auto_ban_threshold_default' THEN '10'
  WHEN 'security_attack_auto_ban_threshold_sys_error_rate_limit' THEN '30'
  WHEN 'security_attack_auto_ban_threshold_scan_upload_ticket_rate_limit' THEN '30'
  WHEN 'security_attack_auto_ban_threshold_scan_upload_ticket_invalid' THEN '50'
  WHEN 'security_sys_error_create_rate_limit_per_minute' THEN '30'
  WHEN 'security_sys_error_create_rate_limit_window_seconds' THEN '60'
  WHEN 'security_visitor_heartbeat_rate_limit_per_minute' THEN '120'
  WHEN 'security_visitor_heartbeat_dedupe_seconds' THEN '3'
  WHEN 'security_tryon_create_rate_limit_per_minute' THEN '20'
  WHEN 'security_tryon_create_concurrency_limit' THEN '2'
  WHEN 'security_cleanup_legacy_casbin_overgrant' THEN 'false'
  ELSE config_value
END
WHERE config_group = 'security'
  AND config_key IN (
    'security_upload_max_size_mb',
    'security_upload_strict_validation_enabled',
    'security_scan_upload_ticket_ttl_seconds',
    'security_scan_upload_ticket_upload_rate_limit_per_ip',
    'security_scan_upload_ticket_upload_rate_limit_window_seconds',
    'security_export_allow_query_token',
    'security_ws_allow_query_token',
    'security_ws_max_conns_per_ip',
    'captcha_rate_limit',
    'captcha_rate_limit_window_seconds',
    'security_login_ip_rate_limit_per_minute',
    'security_login_ip_rate_limit_window_seconds',
    'security_init_api_enabled',
    'security_init_api_private_network_only',
    'security_trusted_proxies',
    'security_public_rate_limit_enabled',
    'security_public_rate_limit_window_seconds',
    'security_public_rate_limit_max_requests',
    'security_attack_auto_ban_enabled',
    'security_attack_auto_ban_window_seconds',
    'security_attack_auto_ban_duration_minutes',
    'security_attack_auto_ban_threshold_default',
    'security_attack_auto_ban_threshold_sys_error_rate_limit',
    'security_attack_auto_ban_threshold_scan_upload_ticket_rate_limit',
    'security_attack_auto_ban_threshold_scan_upload_ticket_invalid',
    'security_sys_error_create_rate_limit_per_minute',
    'security_sys_error_create_rate_limit_window_seconds',
    'security_visitor_heartbeat_rate_limit_per_minute',
    'security_visitor_heartbeat_dedupe_seconds',
    'security_tryon_create_rate_limit_per_minute',
    'security_tryon_create_concurrency_limit',
    'security_cleanup_legacy_casbin_overgrant'
  );
```

After changing values, restart backend service to ensure all runtime paths use new settings.

## 1. Pre-check And Backup

Run a backup before enabling any cleanup.

### MySQL Example

```sql
CREATE TABLE IF NOT EXISTS casbin_rule_backup_20260518 LIKE casbin_rule;
INSERT INTO casbin_rule_backup_20260518 SELECT * FROM casbin_rule;
```

If you run this multiple times, use a unique suffix per run (for example a timestamp).

## 2. One-Time Casbin Legacy Cleanup

The cleanup logic is safe by default because it is disabled unless explicitly enabled.

### Windows PowerShell (temporary for current shell)

```powershell
$env:CS_CLEANUP_LEGACY_CASBIN_OVERGRANT = "true"
```

Restart the server once, wait for successful startup, then disable it:

```powershell
$env:CS_CLEANUP_LEGACY_CASBIN_OVERGRANT = "false"
```

### Linux Shell (temporary for current shell)

```bash
export CS_CLEANUP_LEGACY_CASBIN_OVERGRANT=true
```

Restart once, then disable:

```bash
export CS_CLEANUP_LEGACY_CASBIN_OVERGRANT=false
```

## 3. Verification Queries

### 3.1 Sensitive routes should only remain for admin roles (888, 8881)

```sql
SELECT v0, v1, v2
FROM casbin_rule
WHERE ptype = 'p'
  AND (v1, v2) IN (
    ('/order/deleteOrder', 'DELETE'),
    ('/order/deleteOrderByIds', 'DELETE'),
    ('/order/confirmPayment', 'POST'),
    ('/order/refundOrder', 'POST'),
    ('/order/batchUpdateOrderStatus', 'POST'),
    ('/order/findOrder', 'GET'),
    ('/order/getOrderList', 'GET'),
    ('/sysConfig/updateSysConfig', 'PUT'),
    ('/sysConfig/getSysConfigList', 'GET'),
    ('/sysConfig/getConfigByKey', 'GET'),
    ('/sysConfig/getSysConfigByGroup', 'GET'),
    ('/sysConfig/getSysConfigByKey', 'GET'),
    ('/sysConfig/getAliyunTryonQuotaEstimate', 'GET'),
    ('/sysConfig/getModelCallLogList', 'GET')
  )
  AND v0 NOT IN ('888', '8881');
```

Expected result: zero rows.

### 3.2 Module-managed routes should not keep unknown roles

```sql
SELECT DISTINCT v0
FROM casbin_rule
WHERE ptype = 'p'
  AND v1 IN (
    '/order/createOrder',
    '/order/placeOrder',
    '/sysConfig/getPaymentConfig',
    '/sysConfig/getUniPreferredPayConfig',
    '/tryonTask/createTryonTask'
  )
  AND v0 NOT IN ('888', '8881', '8080', '9528');
```

Expected result: zero rows.

### 3.3 Hotlink folder listing should not be granted to role 8080

```sql
SELECT v0, v1, v2
FROM casbin_rule
WHERE ptype = 'p'
  AND v1 = '/fileUploadAndDownload/listFolders'
  AND v2 = 'GET'
  AND v0 NOT IN ('888', '8881', '9528');
```

Expected result: zero rows.

### 3.4 File ownership rows should be trackable

```sql
SELECT
  SUM(CASE WHEN created_by = 0 THEN 1 ELSE 0 END) AS no_owner_rows,
  COUNT(*) AS total_rows
FROM exa_file_upload_and_downloads;
```

If `no_owner_rows` is high, prefer using role `888` for cleanup/backfill before assigning daily operations to non-super file managers.

## 4. API Smoke Checks

After deployment:
- Role 8080 must fail on admin-only endpoints, such as `/order/confirmPayment` and `/sysConfig/getSysConfigByGroup`.
- Role 8080 must fail on file library management endpoints, including `/fileUploadAndDownload/deleteFile`, `/fileUploadAndDownload/getFileList`, and `/fileUploadAndDownload/listFolders`.
- Roles 888/8881 must still pass on admin endpoints.
- Public `getTryonConfig` should not expose token/secret fields inside `tryon_models` payload.
- QR scan upload URL should contain `ticket=` and must not contain `token=`.
- Reusing the same scan upload ticket for a second upload should fail.
- Repeated high-frequency calls to `uploadByTicket` from same `IP + User-Agent` should hit `requestTooFrequent`.
- Trigger repeated invalid/expired ticket requests and verify `/sysBannedIP/getAttackStats` shows `scanTicketInvalid` growth.
- Trigger repeated rate-limited uploadByTicket requests and verify `/sysBannedIP/getAttackStats` shows `scanTicketRate` growth.

## 5. Export Token Transport Rollout

The default behavior is now safer:
- Download token via request header is supported.
- Query token can be blocked by default.

If old clients still rely on URL query token for export download, temporarily enable:

```bash
export CS_EXPORT_ALLOW_QUERY_TOKEN=true
```

After clients are upgraded, set back to `false`.

## 6. Upload Guardrail Rollout

Default max upload size is 20 MB.

To override:

```bash
export CS_UPLOAD_MAX_SIZE_MB=20
```

Security checks now reject:
- empty files
- dangerous executable/script extensions
- blocked dangerous MIME/content types (sniffed from file bytes)
- invalid filenames

## 7. Rollback Plan

If cleanup produced unexpected access loss, restore from backup:

```sql
TRUNCATE TABLE casbin_rule;
INSERT INTO casbin_rule SELECT * FROM casbin_rule_backup_20260518;
```

Then restart the server and keep cleanup switch disabled:

```bash
export CS_CLEANUP_LEGACY_CASBIN_OVERGRANT=false
```

## 8. Recommended Rollout Order

1. Backup `casbin_rule`.
2. Deploy code.
3. Enable cleanup switch for one startup window.
4. Run verification SQL.
5. Disable cleanup switch.
6. Keep `CS_EXPORT_ALLOW_QUERY_TOKEN` only as temporary compatibility fallback.

## 9. Uni Dependency Remediation (Manual, Rollback-Friendly)

Current verified baseline:
- `npm audit` on `uni/`: `51` vulnerabilities (`34 high`, `8 moderate`, `9 low`).
- `npm run build:h5`: exit code `0`.

Important constraints validated in this workspace:
- Single-package bump `vite -> 5.4.21` increased vulnerabilities (`51 -> 54`) under current `@dcloudio/uni-* 3.0.0-405...` ecosystem.
- `npm audit fix --force` is currently unreliable and fails with:
  - `ETARGET No matching version found for undefined@undefined`
- Full-train DCloud alpha upgrade (`@dcloudio/* -> 3.0.0-alpha-5000920260515001`) keeps build green but increases vulnerabilities (`51 -> 61`), so it is rejected by release gate and should be rolled back.

Do not use one-shot force remediation in CI. Use grouped manual upgrades with checkpoint validation.

### 9.1 Upgrade Grouping

Group A (highest coupling, must upgrade together):
- `@dcloudio/uni-app`
- `@dcloudio/uni-app-plus`
- `@dcloudio/uni-h5`
- `@dcloudio/vite-plugin-uni`
- `@dcloudio/uni-cli-shared`

Group B (platform adapters, upgrade after Group A is stable):
- `@dcloudio/uni-mp-*`
- `@dcloudio/uni-quickapp-webview`
- `@dcloudio/uni-app-harmony`

Group C (transitive chain usually fixed by Group A):
- `@intlify/core-base`
- `@intlify/message-resolver`
- `esbuild`

### 9.2 One Group = One Change Window

For each group:

1. Backup manifests.

```powershell
Copy-Item package.json package.json.group.bak -Force
Copy-Item package-lock.json package-lock.json.group.bak -Force
```

2. Upgrade only one group (pin exact versions).

3. Validate immediately.

```powershell
npm install
npm run build:h5
npm audit --omit=dev
npm audit
```

4. If build or runtime smoke fails, rollback immediately.

```powershell
Copy-Item package.json.group.bak package.json -Force
Copy-Item package-lock.json.group.bak package-lock.json -Force
npm install
```

### 9.3 Release Gate Recommendation

Before shipping:
- Build must pass (`npm run build:h5` exit `0`).
- Vulnerability count must be non-increasing versus baseline.
- Any remaining high vulnerabilities must have explicit risk acceptance and owner.

### 9.4 Risk Acceptance Tracking

Use `docs/uni-vulnerability-risk-register.md` as the source of truth for:
- temporary risk acceptance
- owner assignment
- target remediation date

### 9.5 Automated Gate Command

Use the repository command below in `uni/`:

```bash
npm run security:gate
```

CI integration:
- `.github/workflows/ci.yaml` (`frontend-uni` job) now runs `npm run security:gate` with baseline env values.

The command runs:
1. `npm run build:h5`
2. `npm audit --json`
3. baseline comparison (`total/high/moderate/low` must be non-increasing)

Default baseline values are:
- total: `51`
- high: `34`
- moderate: `8`
- low: `9`

Override when needed:

```bash
UNI_AUDIT_BASELINE_TOTAL=51 UNI_AUDIT_BASELINE_HIGH=34 UNI_AUDIT_BASELINE_MODERATE=8 UNI_AUDIT_BASELINE_LOW=9 npm run security:gate
```
