//go:build integration
// +build integration

package spanner

import (
	"context"
	"testing"
	"time"
)

func TestPermissionRepository_GetPermissionsForAccountGuid(t *testing.T) {
	repository := connectToPermissionRepository()
	tests := []struct {
		name string
		arg  string
	}{
		{
			name: "get_permissions_by_person_guuid",
			arg:  "2c6984f3-af80-4bfd-81c7-0f597c12662e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms, err := repository.GetPermissionsForAccountGuid(context.Background(), tt.arg)
			if err != nil {
				t.Errorf("got: %v", err)
			}
			if tt.arg != "" && len(perms) < 0 {
				t.Errorf("Expected permissions, got: %v", len(perms))
			}
		})
	}
}
func TestPermissionRepository_GetPermissionsForAPIKey(t *testing.T) {
	repository := connectToPermissionRepository()
	tests := []struct {
		name string
		arg  string
	}{
		{
			name: "get_create_permission",
			arg:  "6247c10c-84a0-4fa1-b330-77eea1e944d3",
		},
		{
			name: "get_permission_with_non_existing_api_key",
			arg:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms, err := repository.GetPermissionsForAPIKey(context.Background(), tt.arg)
			if err != nil {
				t.Errorf("got: %v", err)
			}
			if tt.arg != "" && len(perms) < 0 {
				t.Errorf("Expected permissions, got: %v", len(perms))
			}
		})
	}
}

func connectToPermissionRepository() *PermissionRepository {
	wtimeout, _ := time.ParseDuration("10s")
	rtimeout, _ := time.ParseDuration("10s")
	return NewPermissionRepository(Client, rtimeout, wtimeout)
}
