package system

import "testing"

func TestSysParamsManageRoleRequiresSuperAdmin(t *testing.T) {
	if !isSysParamsManageRole(888) {
		t.Fatal("expected 888 to manage system params")
	}
	if isSysParamsManageRole(8881) {
		t.Fatal("expected 8881 to be denied system params management")
	}
}
