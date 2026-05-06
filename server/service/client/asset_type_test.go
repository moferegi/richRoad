package client

import (
	"testing"

	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
)

func strPtr(s string) *string {
	return &s
}

func TestNormalizeAssetType(t *testing.T) {
	tests := []struct {
		name    string
		input   *string
		want    string
		wantErr bool
	}{
		{name: "nil defaults to point", input: nil, want: clientModel.AssetTypePoint},
		{name: "empty defaults to point", input: strPtr("   "), want: clientModel.AssetTypePoint},
		{name: "point remains point", input: strPtr(clientModel.AssetTypePoint), want: clientModel.AssetTypePoint},
		{name: "tryon remains tryon", input: strPtr(clientModel.AssetTypeTryonPoint), want: clientModel.AssetTypeTryonPoint},
		{name: "invalid asset type", input: strPtr("coin"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAssetType(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestBuildPointRecord_DefaultAssetType(t *testing.T) {
	record := buildPointRecord(7, "", "increase", 3, "guest_init", "游客赠送", "remark")
	if record.AssetType == nil {
		t.Fatalf("asset type should not be nil")
	}
	if *record.AssetType != clientModel.AssetTypePoint {
		t.Fatalf("expected default asset type %s, got %s", clientModel.AssetTypePoint, *record.AssetType)
	}
	if record.UserId == nil || *record.UserId != 7 {
		t.Fatalf("unexpected user id: %+v", record.UserId)
	}
}

func TestBuildPointRecord_ExplicitAssetType(t *testing.T) {
	record := buildPointRecord(8, clientModel.AssetTypeTryonPoint, "decrease", 1, "consume", "试衣扣费", "remark")
	if record.AssetType == nil {
		t.Fatalf("asset type should not be nil")
	}
	if *record.AssetType != clientModel.AssetTypeTryonPoint {
		t.Fatalf("expected asset type %s, got %s", clientModel.AssetTypeTryonPoint, *record.AssetType)
	}
	if record.PointChange == nil || *record.PointChange != 1 {
		t.Fatalf("unexpected point change: %+v", record.PointChange)
	}
}
