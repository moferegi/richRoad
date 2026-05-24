package system

import (
	"testing"

	systemModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func TestApiTokenManageRoleRequiresSuperAdmin(t *testing.T) {
	if !isApiTokenManageRole(888) {
		t.Fatal("expected 888 to manage API tokens")
	}
	if isApiTokenManageRole(8881) {
		t.Fatal("expected 8881 to be denied API token management")
	}
}

func TestRedactApiTokenList(t *testing.T) {
	list := []systemModel.SysApiToken{
		{Token: "header.payload.signature", Remark: "active token"},
		{Token: "", Remark: "empty token"},
	}

	redacted := redactApiTokenList(list)
	if redacted[0].Token != apiTokenListPlaceholder {
		t.Fatalf("expected stored token to be redacted, got %q", redacted[0].Token)
	}
	if redacted[1].Token != "" {
		t.Fatalf("expected empty token to stay empty, got %q", redacted[1].Token)
	}
	if list[0].Token != "header.payload.signature" {
		t.Fatalf("redaction mutated original list token: %q", list[0].Token)
	}
}
