package utils

import (
	"testing"
)

func TestGetJSONKeys(t *testing.T) {
	var jsonStr = `
	{
		"Name": "test",
		"TableName": "test",
		"TemplateID": "test",
		"TemplateInfo": "test",
		"Limit": 0
}`
	keys, err := GetJSONKeys(jsonStr)
	if err != nil {
		t.Fatalf("GetJSONKeys failed: %v", err)
	}
	if len(keys) != 5 {
		t.Fatalf("GetJSONKeys length mismatch: got %d, want 5", len(keys))
	}
	if keys[0] != "Name" {
		t.Fatalf("GetJSONKeys index 0 mismatch: got %q, want %q", keys[0], "Name")
	}
	if keys[1] != "TableName" {
		t.Fatalf("GetJSONKeys index 1 mismatch: got %q, want %q", keys[1], "TableName")
	}
	if keys[2] != "TemplateID" {
		t.Fatalf("GetJSONKeys index 2 mismatch: got %q, want %q", keys[2], "TemplateID")
	}
	if keys[3] != "TemplateInfo" {
		t.Fatalf("GetJSONKeys index 3 mismatch: got %q, want %q", keys[3], "TemplateInfo")
	}
	if keys[4] != "Limit" {
		t.Fatalf("GetJSONKeys index 4 mismatch: got %q, want %q", keys[4], "Limit")
	}
}
