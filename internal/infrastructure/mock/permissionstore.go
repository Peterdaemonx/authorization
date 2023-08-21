package mock

import (
	"context"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

var permissions = []processing.Permission{
	{
		ID:     uuid.New().String(),
		Code:   "create_authorization",
		Label:  "Create new authorizations",
		ApiKey: "6247c10c-84a0-4fa1-b330-77eea1e944d3",
	},
	{
		ID:     uuid.New().String(),
		Code:   "create_refund",
		Label:  "Create new refund",
		ApiKey: "6247c10c-84a0-4fa1-b330-77eea1e944d3",
	},
	{
		ID:     uuid.New().String(),
		Code:   "create_reversal",
		Label:  "Create new reversal",
		ApiKey: "6247c10c-84a0-4fa1-b330-77eea1e944d3",
	},
	{
		ID:     uuid.New().String(),
		Code:   "create_capture",
		Label:  "Create new capture",
		ApiKey: "6247c10c-84a0-4fa1-b330-77eea1e944d3",
	},
	{
		ID:     uuid.New().String(),
		Code:   "get_authorizations",
		Label:  "Get all authorizations",
		ApiKey: "6247c10c-84a0-4fa1-b330-77eea1e944d3",
	},
	{
		ID:     uuid.New().String(),
		Code:   "get_captures",
		Label:  "Get status of reversals and captures",
		ApiKey: "4c11c341-af5f-487d-860e-7659fad6288c",
	},
}

type PermissionStore struct {
}

func (p PermissionStore) GetPspForAPIKey(ctx context.Context, apiKey string) (entity.PSP, error) {
	return entity.PSP{}, nil
}

func (p PermissionStore) GetPermissionsForAPIKey(_ context.Context, apiKey string) ([]processing.Permission, error) {
	return permissionsFromApiKey(apiKey)
}

func permissionsFromApiKey(apiKey string) ([]processing.Permission, error) {
	var permissionsFound []processing.Permission
	for _, permission := range permissions {
		if permission.ApiKey == apiKey {
			permissionsFound = append(permissionsFound, permission)
		}
	}
	return permissionsFound, nil
}

func (p PermissionStore) GetPermissionsForAccountGuid(_ context.Context, _ string) ([]processing.Permission, error) {
	return permissions, nil
}
