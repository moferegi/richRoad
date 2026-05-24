package system

import "testing"

func TestApiManageRoleRequiresSuperAdmin(t *testing.T) {
	if !isApiManageRole(888) {
		t.Fatal("expected 888 to manage API definitions")
	}
	if isApiManageRole(8881) {
		t.Fatal("expected 8881 to be denied API definition management")
	}
}

func TestCasbinManageRoleRequiresSuperAdmin(t *testing.T) {
	if !isCasbinManageRole(888) {
		t.Fatal("expected 888 to manage Casbin policies")
	}
	if isCasbinManageRole(8881) {
		t.Fatal("expected 8881 to be denied Casbin policy management")
	}
}

func TestUserAuthoritySwitchRoleGuard(t *testing.T) {
	if isUserManageRole(9528) {
		t.Fatal("expected 9528 to be denied system user management")
	}
	if !canManageUserAuthority(888, 888) {
		t.Fatal("expected 888 to manage super admin authority")
	}
	if canManageUserAuthority(8881, 888) {
		t.Fatal("expected 8881 to be denied switching users to 888")
	}
}

func TestAuthorityButtonManageRoleGuard(t *testing.T) {
	if !isAuthorityBtnManageRole(888) {
		t.Fatal("expected 888 to manage authority buttons")
	}
	if !isAuthorityBtnManageRole(8881) {
		t.Fatal("expected 8881 to manage scoped authority buttons")
	}
	if isAuthorityBtnManageRole(9528) {
		t.Fatal("expected 9528 to be denied authority button management")
	}
}

func TestAuthorityCopyTargetGuard(t *testing.T) {
	if !canManageAuthorityTarget(888, 888) {
		t.Fatal("expected 888 to copy/manage 888 authority")
	}
	if canManageAuthorityTarget(8881, 888) {
		t.Fatal("expected 8881 to be denied copying/managing 888 authority")
	}
	if !canManageAuthorityTarget(8881, 9528) {
		t.Fatal("expected 8881 to manage non-888 authority targets")
	}
}
