package system

import (
	"testing"

	systemModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func buildSafeExportTemplate() systemModel.SysExportTemplate {
	limit := 100
	return systemModel.SysExportTemplate{
		Name:         "safe-template",
		TableName:    "sys_user",
		TemplateID:   "tpl_safe",
		TemplateInfo: `{"id":"ID","username":"用户名"}`,
		Limit:        &limit,
		Order:        "id desc",
		Conditions: []systemModel.Condition{
			{From: "username", Column: "sys_user.username", Operator: "LIKE"},
		},
		JoinTemplate: []systemModel.JoinTemplate{
			{JOINS: "LEFT JOIN", Table: "sys_authority", ON: "sys_user.authority_id = sys_authority.authority_id"},
		},
	}
}

func TestValidateExportTemplateQuerySafety_AllowsSafeTemplate(t *testing.T) {
	tpl := buildSafeExportTemplate()
	if err := validateExportTemplateQuerySafety(tpl); err != nil {
		t.Fatalf("expected safe template to pass, got error: %v", err)
	}
}

func TestValidateExportTemplateQuerySafety_RejectsInvalidOperator(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.Conditions[0].Operator = "OR"

	err := validateExportTemplateQuerySafety(tpl)
	if err == nil || err.Error() != errExportTemplateConditionOperatorInvalid {
		t.Fatalf("expected operator validation error, got: %v", err)
	}
}

func TestValidateExportTemplateQuerySafety_RejectsInvalidColumn(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.Conditions[0].Column = "sys_user.username;drop table sys_user"

	err := validateExportTemplateQuerySafety(tpl)
	if err == nil || err.Error() != errExportTemplateConditionColumnInvalid {
		t.Fatalf("expected column validation error, got: %v", err)
	}
}

func TestValidateExportTemplateQuerySafety_RejectsRawSQLByDefault(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.SQL = "SELECT id, username FROM sys_user WHERE username = @username"
	t.Setenv(exportAllowRawSQLEnvKey, "false")

	err := validateExportTemplateQuerySafety(tpl)
	if err == nil || err.Error() != errExportTemplateRawSQLDisabled {
		t.Fatalf("expected raw sql to be blocked by default, got: %v", err)
	}
}

func TestValidateExportTemplateQuerySafety_AllowsSafeRawSQLWhenEnabled(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.SQL = "SELECT id, username FROM sys_user WHERE username = @username"
	t.Setenv(exportAllowRawSQLEnvKey, "true")

	if err := validateExportTemplateQuerySafety(tpl); err != nil {
		t.Fatalf("expected safe raw sql to pass when enabled, got: %v", err)
	}
}

func TestValidateExportTemplateQuerySafety_RejectsDangerousRawSQL(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.SQL = "SELECT id FROM sys_user; UPDATE sys_user SET username='x'"
	t.Setenv(exportAllowRawSQLEnvKey, "true")

	err := validateExportTemplateQuerySafety(tpl)
	if err == nil {
		t.Fatalf("expected dangerous raw sql to be rejected")
	}
}

func TestValidateExportTemplateQuerySafety_RejectsImportSQLByDefault(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.ImportSQL = "INSERT INTO sys_user(username) VALUES(@username)"
	t.Setenv(exportAllowImportSQLEnvKey, "false")

	err := validateExportTemplateQuerySafety(tpl)
	if err == nil || err.Error() != errExportTemplateImportSQLDisabled {
		t.Fatalf("expected import sql to be blocked by default, got: %v", err)
	}
}

func TestValidateExportTemplateQuerySafety_AllowsSafeImportSQLWhenEnabled(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.ImportSQL = "INSERT INTO sys_user(username) VALUES(@username)"
	t.Setenv(exportAllowImportSQLEnvKey, "true")

	if err := validateExportTemplateQuerySafety(tpl); err != nil {
		t.Fatalf("expected safe import sql to pass when enabled, got: %v", err)
	}
}

func TestValidateExportTemplateQuerySafety_RejectsDangerousImportSQL(t *testing.T) {
	tpl := buildSafeExportTemplate()
	tpl.ImportSQL = "UPDATE sys_user SET username=@username"
	t.Setenv(exportAllowImportSQLEnvKey, "true")

	err := validateExportTemplateQuerySafety(tpl)
	if err == nil || err.Error() != errExportTemplateImportSQLWhereRequired {
		t.Fatalf("expected dangerous import sql to be rejected, got: %v", err)
	}
}
