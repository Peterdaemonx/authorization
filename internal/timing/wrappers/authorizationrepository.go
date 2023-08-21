// Code generated by timing/wrappers/generate, DO NOT EDIT.
// Generated on Mon May 15 15:07 2023
package timingwrappers

import (
	"context"
	uuid "github.com/google/uuid"
	app "gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app"
	entity "gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	timing "gitlab.cmpayments.local/creditcard/authorization/internal/timing"
)

type AuthorizationRepository struct {
	Base app.Repository
}

func (w AuthorizationRepository) AuthorizationAlreadyReversed(ctx context.Context, id uuid.UUID) (bool, error) {
	timing.Start(ctx, "AuthorizationRepository.AuthorizationAlreadyReversed")
	defer timing.Stop(ctx, "AuthorizationRepository.AuthorizationAlreadyReversed")
	return w.Base.AuthorizationAlreadyReversed(ctx, id)
}
func (w AuthorizationRepository) CreateAuthorization(ctx context.Context, a entity.Authorization) error {
	timing.Start(ctx, "AuthorizationRepository.CreateAuthorization")
	defer timing.Stop(ctx, "AuthorizationRepository.CreateAuthorization")
	return w.Base.CreateAuthorization(ctx, a)
}
func (w AuthorizationRepository) CreateMastercardAuthorization(ctx context.Context, a entity.Authorization) error {
	timing.Start(ctx, "AuthorizationRepository.CreateMastercardAuthorization")
	defer timing.Stop(ctx, "AuthorizationRepository.CreateMastercardAuthorization")
	return w.Base.CreateMastercardAuthorization(ctx, a)
}
func (w AuthorizationRepository) CreateVisaAuthorization(ctx context.Context, a entity.Authorization) error {
	timing.Start(ctx, "AuthorizationRepository.CreateVisaAuthorization")
	defer timing.Stop(ctx, "AuthorizationRepository.CreateVisaAuthorization")
	return w.Base.CreateVisaAuthorization(ctx, a)
}
func (w AuthorizationRepository) GetAllAuthorizations(ctx context.Context, pspID uuid.UUID, filters entity.Filters, params map[string]interface{}) (entity.Metadata, []entity.Authorization, error) {
	timing.Start(ctx, "AuthorizationRepository.GetAllAuthorizations")
	defer timing.Stop(ctx, "AuthorizationRepository.GetAllAuthorizations")
	return w.Base.GetAllAuthorizations(ctx, pspID, filters, params)
}
func (w AuthorizationRepository) GetAuthorization(ctx context.Context, pspID uuid.UUID, authorizationID uuid.UUID) (entity.Authorization, error) {
	timing.Start(ctx, "AuthorizationRepository.GetAuthorization")
	defer timing.Stop(ctx, "AuthorizationRepository.GetAuthorization")
	return w.Base.GetAuthorization(ctx, pspID, authorizationID)
}
func (w AuthorizationRepository) GetAuthorizationWithSchemeData(ctx context.Context, pspID uuid.UUID, authorizationID uuid.UUID) (entity.Authorization, error) {
	timing.Start(ctx, "AuthorizationRepository.GetAuthorizationWithSchemeData")
	defer timing.Stop(ctx, "AuthorizationRepository.GetAuthorizationWithSchemeData")
	return w.Base.GetAuthorizationWithSchemeData(ctx, pspID, authorizationID)
}
func (w AuthorizationRepository) UpdateAuthorizationResponse(ctx context.Context, a entity.Authorization) error {
	timing.Start(ctx, "AuthorizationRepository.UpdateAuthorizationResponse")
	defer timing.Stop(ctx, "AuthorizationRepository.UpdateAuthorizationResponse")
	return w.Base.UpdateAuthorizationResponse(ctx, a)
}
func (w AuthorizationRepository) UpdateAuthorizationStatus(ctx context.Context, authorizationID uuid.UUID, status entity.Status) error {
	timing.Start(ctx, "AuthorizationRepository.UpdateAuthorizationStatus")
	defer timing.Stop(ctx, "AuthorizationRepository.UpdateAuthorizationStatus")
	return w.Base.UpdateAuthorizationStatus(ctx, authorizationID, status)
}