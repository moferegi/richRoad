# Security Hardening Runbook

This runbook documents rollout, verification, and rollback steps for the recent security hardening changes.

## Scope

The runbook covers:
- Casbin legacy over-grant cleanup switch: `CS_CLEANUP_LEGACY_CASBIN_OVERGRANT`
- Export download token transport switch: `CS_EXPORT_ALLOW_QUERY_TOKEN`
- Upload size guardrail: `CS_UPLOAD_MAX_SIZE_MB`

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

## Runtime Config Priority

Security runtime settings now support backend management.

Priority order:
1. `client_sys_config.config_value` (admin panel)
2. Environment variable fallback
3. Built-in default value

### Backend Keys (group: `security`)

- `security_upload_max_size_mb`
- `security_upload_strict_validation_enabled`
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
