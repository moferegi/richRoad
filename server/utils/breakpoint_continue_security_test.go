package utils

import "testing"

func TestIsSafePathSegment(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{name: "simple", in: "file123", want: true},
		{name: "dotdot", in: "../x", want: false},
		{name: "slash", in: "a/b", want: false},
		{name: "backslash", in: `a\\b`, want: false},
		{name: "empty", in: "", want: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isSafePathSegment(c.in); got != c.want {
				t.Fatalf("isSafePathSegment(%q)=%v want=%v", c.in, got, c.want)
			}
		})
	}
}

func TestBreakPointContinue_RejectsUnsafeName(t *testing.T) {
	if _, err := BreakPointContinue([]byte("x"), "../evil", 1, 1, "md5"); err == nil {
		t.Fatalf("expected error for unsafe fileName")
	}
	if _, err := BreakPointContinue([]byte("x"), "ok", 1, 1, "../md5"); err == nil {
		t.Fatalf("expected error for unsafe fileMd5")
	}
}

func TestRemoveChunk_RejectsUnsafePath(t *testing.T) {
	if err := RemoveChunk("../bad"); err == nil {
		t.Fatalf("expected error for unsafe remove path")
	}
}
