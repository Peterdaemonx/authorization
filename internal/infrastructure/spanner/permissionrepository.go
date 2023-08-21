package spanner

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
	"google.golang.org/api/iterator"
)

type PermissionRepository struct {
	client       *spanner.Client
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewPermissionRepository(
	client *spanner.Client,
	readTimeout time.Duration,
	writeTimeout time.Duration) *PermissionRepository {
	return &PermissionRepository{
		client:       client,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (p PermissionRepository) GetPermissionsForAPIKey(ctx context.Context, apiKey string) ([]processing.Permission, error) {
	stmt := spanner.Statement{
		SQL: `
		SELECT permissions.permission_id, permissions.code, permissions.label, ac.api_key
        FROM permissions
        JOIN api_consumers_permissions acp ON permissions.permission_id = acp.permission_id
        JOIN api_consumers ac on acp.api_consumer_id = ac.api_consumer_id
        WHERE ac.api_key = @api_key
	`,
		Params: map[string]interface{}{"api_key": apiKey},
	}
	var permissions []processing.Permission

	ctx, cancel := context.WithTimeout(ctx, p.readTimeout)
	defer cancel()

	iter := p.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return permissions, nil
		}
		if err != nil {
			return nil, fmt.Errorf("cannot iterate over rows: %w", err)
		}
		var permRecord permissionRecord
		if err := row.ToStruct(&permRecord); err != nil {
			return nil, fmt.Errorf("cannot process permission to struct: %w", err)
		}
		permissions = append(permissions, mapPermissionRecordToEntity(permRecord))
	}
}

type permissionRecord struct {
	ID     string `spanner:"permission_id"`
	Code   string `spanner:"code"`
	Label  string `spanner:"label"`
	ApiKey string `spanner:"api_key"`
}

func mapPermissionRecordToEntity(record permissionRecord) processing.Permission {
	return processing.Permission{
		ID:     record.ID,
		Code:   record.Code,
		Label:  record.Label,
		ApiKey: record.ApiKey,
	}
}

func (p PermissionRepository) GetPermissionsForAccountGuid(ctx context.Context, accountGuid string) ([]processing.Permission, error) {
	stmt := spanner.Statement{
		SQL: `SELECT u.person_guid, u.username, permissions.code, permissions.label
				 FROM permissions
				 JOIN users_permissions up ON permissions.permission_id = up.permission_id
				 JOIN users u on u.person_guid = up.person_guid
				 WHERE up.person_guid = @personGuid
				`,
		Params: map[string]interface{}{"personGuid": accountGuid},
	}
	var permissions []processing.Permission

	ctx, cancel := context.WithTimeout(ctx, p.readTimeout)
	defer cancel()

	iter := p.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return permissions, nil
		}
		if err != nil {
			return nil, fmt.Errorf("cannot iterate over rows: %w", err)
		}
		var accountPermRecord accountPermissions
		if err := row.ToStruct(&accountPermRecord); err != nil {
			return nil, fmt.Errorf("cannot process permission to struct: %w", err)
		}
		permissions = append(permissions, mapAccountPermissionRecordToEntity(accountPermRecord))
	}
}

type accountPermissions struct {
	ID       string `spanner:"person_guid"`
	UserName string `spanner:"username"`
	Code     string `spanner:"code"`
	Label    string `spanner:"label"`
}

func mapAccountPermissionRecordToEntity(record accountPermissions) processing.Permission {
	return processing.Permission{
		ID:    record.ID,
		Code:  record.Code,
		Label: record.Label,
	}
}
