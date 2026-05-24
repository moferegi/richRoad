package client

import (
	"strings"
	"testing"

	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
)

func TestRedactSysConfigListForNonSuperAdmin(t *testing.T) {
	models := `[{"key":"main","providerToken":"dashscope-token","beautifyAccessKeySecret":"aliyun-secret","model":"aitryon-plus"}]`
	list := []clientModel.SysConfig{
		{ConfigKey: "tryon_models", ConfigValue: models},
		{ConfigKey: "aws_secret_key", ConfigValue: "cloud-secret"},
		{ConfigKey: "security_ws_allow_query_token", ConfigValue: "false"},
		{ConfigKey: "tryon_cost_points", ConfigValue: "1"},
	}

	redacted := redactSysConfigListForAuthority(list, 8881)

	if strings.Contains(redacted[0].ConfigValue, "dashscope-token") || strings.Contains(redacted[0].ConfigValue, "aliyun-secret") {
		t.Fatalf("tryon_models leaked sensitive fields after redaction: %s", redacted[0].ConfigValue)
	}
	if !strings.Contains(redacted[0].ConfigValue, "aitryon-plus") {
		t.Fatalf("tryon_models redaction removed non-sensitive model metadata: %s", redacted[0].ConfigValue)
	}
	if redacted[1].ConfigValue != sysConfigSecretPlaceholder {
		t.Fatalf("expected direct secret placeholder, got %q", redacted[1].ConfigValue)
	}
	if redacted[2].ConfigValue != "false" {
		t.Fatalf("expected non-sensitive query-token switch to remain visible, got %q", redacted[2].ConfigValue)
	}
	if redacted[3].ConfigValue != "1" {
		t.Fatalf("expected normal config to remain visible, got %q", redacted[3].ConfigValue)
	}
}

func TestRedactSysConfigListKeepsSuperAdminValues(t *testing.T) {
	list := []clientModel.SysConfig{{ConfigKey: "aws_secret_key", ConfigValue: "cloud-secret"}}
	redacted := redactSysConfigListForAuthority(list, 888)
	if redacted[0].ConfigValue != "cloud-secret" {
		t.Fatalf("expected super admin to receive original value, got %q", redacted[0].ConfigValue)
	}
}
