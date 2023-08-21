package ports

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/visa"
	"gitlab.cmpayments.local/creditcard/platform"
	platformErr "gitlab.cmpayments.local/creditcard/platform/http/errors"
	platformhandler "gitlab.cmpayments.local/creditcard/platform/http/handler"
	"gitlab.cmpayments.local/creditcard/platform/http/logging"
	"gitlab.cmpayments.local/creditcard/platform/http/validator"
	liblogging "gitlab.cmpayments.local/libraries-go/logging"

	"gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/tokenization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/cardinfo"
	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

var (
	ErrSubseqSpecifiedWithEmptyTraceId = errors.New("specified as a subsequent authorization but traceId can not be empty")
)

type authorizationHandler struct {
	logger                     platform.Logger
	allowProductionCardNumbers bool
	cardRanges                 *cardinfo.Collection
	authorizationService       app.AuthorizationService
}

func NewAuthorizationHandler(allowProductionCardNumbers bool, cardRanges *cardinfo.Collection, logger platform.Logger, authService app.AuthorizationService) *authorizationHandler {
	return &authorizationHandler{
		allowProductionCardNumbers: allowProductionCardNumbers,
		cardRanges:                 cardRanges,
		logger:                     logger,
		authorizationService:       authService,
	}
}

func (h *authorizationHandler) CreateAuthorization(w http.ResponseWriter, r *http.Request) {
	input := authorizationRequest{}
	ctx := r.Context()

	h.logger.Info(ctx, "create authorization")

	if err := platformhandler.ReadJSON(w, r, &input); err != nil {
		platformErr.BadRequestResponse(ctx, w, h.logger, err)
		return
	}

	cardInfo, ok := h.cardRanges.Find(input.Card.Number)
	if !ok {
		platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"card.number": {"unknown card range"}})
		return
	}

	if cardInfo.IsBlocked {
		platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"card.number": {"card bin is blocked"}})
		return
	}

	input.Card.scheme = cardInfo.Scheme

	recurring, err := mapRecurring(input.InitialTraceID, input.InitialRecurring)
	if err != nil {
		switch {
		case errors.Is(err, ErrSubseqSpecifiedWithEmptyTraceId):
			platformErr.BadRequestResponse(ctx, w, h.logger, err)
			return
		default:
			platformErr.ServerErrorResponse(ctx, w, h.logger, err)
			return
		}
	}

	psp, _ := processing.PSPFromContext(ctx)

	v := validator.New()

	input.validate(v)
	if !v.Valid() {
		platformErr.FailedValidationResponse(ctx, w, h.logger, v.Errors)
		return
	}

	authorization := mapAuthorizationRequest(ctx, cardInfo, psp, input, recurring)

	err = h.authorizationService.Authorize(ctx, &authorization)
	if err != nil {
		switch {
		case errors.Is(err, visa.CavvErrorNumeric):
			h.logger.Error(liblogging.ContextWithError(ctx, err), "unprocessable content")
			platformErr.FailedValidationResponse(ctx, w, h.logger, map[string][]string{"authenticationVerificationValue": {"failed to encode CAVV"}})
			return
		case errors.Is(err, tokenization.ErrFailedTokenize):
			h.logger.Error(liblogging.ContextWithError(ctx, err), "internal server error")
			platformErr.ServerErrorResponse(ctx, w, h.logger, tokenization.ErrFailedTokenize)
			return
		default:
			platformErr.ServerErrorResponse(ctx, w, h.logger, err)
			return
		}
	}

	if err = platformhandler.WriteJSON(w, http.StatusOK, mapAuthorizationResponse(authorization), nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, err)
	}
}

func (h *authorizationHandler) GetAuthorization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	authorizationID := params.ByName("authorization")
	id := uuid.MustParse(authorizationID)
	psp, _ := processing.PSPFromContext(ctx)

	auth, err := h.authorizationService.GetAuthorization(ctx, psp.ID, id)
	if err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, errors.New("unable to fetch authorization"))
		return
	}

	if err := platformhandler.WriteJSON(w, http.StatusOK, mapGetAuthorizationResponse(auth), nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, errors.New("internal server error"))
	}
}

func mapRecurring(traceID string, initialRecurring bool) (entity.Recurring, error) {
	if initialRecurring && traceID == "" {
		return entity.Recurring{
			Initial:    initialRecurring,
			Subsequent: false,
		}, nil
	}

	if !initialRecurring && traceID != "" {
		return entity.Recurring{
			Subsequent: true,
			TraceID:    traceID,
		}, nil
	}

	if !initialRecurring && traceID == "" {
		return entity.Recurring{}, nil
	}

	return entity.Recurring{}, ErrSubseqSpecifiedWithEmptyTraceId
}

func mapAuthorizationRequest(ctx context.Context, cardInfo cardinfo.Range, psp entity.PSP, input authorizationRequest, recurring entity.Recurring) entity.Authorization {
	if input.AuthorizationType == "" {
		input.AuthorizationType = "finalAuthorization"
	}

	authorization := entity.Authorization{
		ID:                       uuid.New(),
		LogID:                    uuid.MustParse(ctx.Value(logging.LogIDKey).(string)),
		Amount:                   input.Amount,
		Currency:                 currencycode.Must(input.Currency),
		CustomerReference:        input.Reference,
		Source:                   entity.Source(input.Source),
		LocalTransactionDateTime: input.LocalTransactionDateTime,
		Recurring:                recurring,
		Card: entity.Card{
			Number:    input.Card.Number,
			MaskedPan: entity.MaskPan(input.Card.Number),
			Cvv:       input.Card.Cvv,
			Holder:    input.Card.Holder,
			Expiry:    entity.Expiry{Year: input.Card.Expiry.Year, Month: input.Card.Expiry.Month},
			Info:      cardInfo,
		},
		CardAcceptor: entity.CardAcceptor{
			CategoryCode: input.CardAcceptor.CategoryCode,
			ID:           input.CardAcceptor.ID,
			Name:         input.CardAcceptor.Name,
			Address: entity.CardAcceptorAddress{
				PostalCode:  input.CardAcceptor.PostalCode,
				City:        input.CardAcceptor.City,
				CountryCode: input.CardAcceptor.Country,
			},
		},
		Psp: entity.PSP{
			ID: psp.ID,
		},
		Exemption: entity.ExemptionType(input.Exemption),
		ThreeDSecure: entity.ThreeDSecure{
			AuthenticationVerificationValue: input.ThreeDSecure.AuthenticationVerificationValue,
			DirectoryServerID:               input.ThreeDSecure.DirectoryServerTransactionID,
			EcommerceIndicator:              int(input.ThreeDSecure.EcommerceIndicator),
			Version:                         input.ThreeDSecure.Version,
		},
		CitMitIndicator: entity.CitMitIndicator{
			InitiatedBy: entity.MapInitiatedByFromStr(input.CitMitIndicator.InitiatedBy),
			SubCategory: entity.MapSubCategoryFromStr(input.CitMitIndicator.SubCategory),
		},
		MastercardSchemeData: entity.MastercardSchemeData{
			Request: entity.MastercardSchemeRequest{
				AuthorizationType: entity.AuthorizationType(input.AuthorizationType),
			},
		},
	}

	return authorization
}

func mapAuthorizationResponse(a entity.Authorization) authorizationResponse {
	var citMitIndicator *CitMitIndicator
	if a.CitMitIndicator != (entity.CitMitIndicator{}) {
		citMitIndicator = &CitMitIndicator{
			InitiatedBy: string(a.CitMitIndicator.InitiatedBy),
			SubCategory: string(a.CitMitIndicator.SubCategory),
		}
	}

	authorization := authorizationResponse{
		ID:                       a.ID.String(),
		LogID:                    a.LogID.String(),
		Amount:                   a.Amount,
		Currency:                 a.Currency.Alpha3(),
		Reference:                a.CustomerReference,
		Source:                   string(a.Source),
		LocalTransactionDateTime: &a.LocalTransactionDateTime,
		AuthorizationType:        string(a.MastercardSchemeData.Request.AuthorizationType),
		InitialRecurring:         a.Recurring.Initial,
		TraceID:                  a.Recurring.TraceID,
		ProcessingDate:           a.ProcessingDate.Format(time.RFC3339),
		Card:                     CardResponse{Number: a.Card.MaskedPan, Scheme: a.Card.Info.Scheme},
		CardAcceptor: CardAcceptor{
			ID:           a.CardAcceptor.ID,
			CategoryCode: a.CardAcceptor.CategoryCode,
			Name:         a.CardAcceptor.Name,
			City:         a.CardAcceptor.Address.City,
			Country:      a.CardAcceptor.Address.CountryCode,
			PostalCode:   a.CardAcceptor.Address.PostalCode,
		},
		Exemption: string(a.Exemption),
		ThreeDSecure: ThreeDSecureResponse{
			Version:                      a.ThreeDSecure.Version,
			EcommerceIndicator:           fmt.Sprintf("%02d", a.CardSchemeData.Response.EcommerceIndicator),
			DirectoryServerTransactionID: a.ThreeDSecure.DirectoryServerID,
		},
		CitMitIndicator: citMitIndicator,
		CardSchemeResponse: CardSchemeResponse{
			Status:  a.CardSchemeData.Response.Status.String(),
			Code:    a.CardSchemeData.Response.ResponseCode.Value,
			Message: a.CardSchemeData.Response.ResponseCode.Description,
			TraceID: a.CardSchemeData.Response.TraceId,
		},
	}

	return authorization
}

func (h *authorizationHandler) GetAuthorizations(w http.ResponseWriter, r *http.Request) {
	var input struct {
		params map[string]interface{}
		entity.Filters
	}

	ctx := r.Context()
	v := validator.New()
	qs := r.URL.Query()

	reference := web.ReadString(qs, "reference", "")
	amount := web.ReadInt(qs, "amount", -1, v)
	processingDate := web.ReadDate(qs, "processingDate", time.Time{}, v)
	pan := web.ReadString(qs, "pan", "")
	status := web.ReadString(qs, "status", "")
	exemption := web.ReadString(qs, "exemption", "")
	responseCode := web.ReadString(qs, "responseCode", "")
	traceID := web.ReadString(qs, "traceID", "")
	startDate := web.ReadDate(qs, "startDate", time.Time{}, v)
	endDate := web.ReadDate(qs, "endDate", time.Time{}, v)

	input.Filters.Page = web.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = web.ReadInt(qs, "pageSize", 15, v)
	input.Filters.Sort = web.ReadString(qs, "sort", "createdAt")
	input.Filters.SortSafelist = []string{"amount", "-amount", "createdAt", "-createdAt", "pspId", "-pspId", "processingDate", "-processingDate", "status", "-status", "exemption", "-exemption"}

	v.Check(len(pan) <= 4, "pan", []string{"cannot be longer than 4 digits"})
	v.Check(len(responseCode) <= 2, "responseCode", []string{"cannot be longer than 2 characters"})
	v.Check(len(traceID) <= 15, "traceID", []string{"cannot be longer than 15 characters"})
	v.Check(v.HasTwoOrNoDates(startDate, endDate), "endDate", []string{"there must be an endDate if a startDate is present"})
	v.Check(v.IsEqualOrAfter(startDate, endDate), "date", []string{"startDate must be before or equal endDate"})

	// No need to validate the PSPID. This is done in the permissionstore middleware
	parameters := map[string]interface{}{
		"reference":      reference,
		"amount":         amount,
		"processingDate": processingDate,
		"pan":            pan,
		"status":         status,
		"exemption":      exemption,
		"responseCode":   responseCode,
		"traceID":        traceID,
		"startDate":      startDate,
		"endDate":        endDate,
	}

	if ValidateFilters(v, input.Filters); !v.Valid() {
		platformErr.FailedValidationResponse(ctx, w, h.logger, v.Errors)
		return
	}

	psp, _ := processing.PSPFromContext(ctx)

	metadata, authorizations, err := h.authorizationService.GetAuthorizations(ctx, psp.ID, input.Filters, parameters)
	if err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, err)
		return
	}

	response := AuthorizationsResponse{
		Authorizations: []authorizationResponse{},
	}

	for _, a := range authorizations {
		response.Authorizations = append(response.Authorizations, mapAuthorizationResponse(a))
	}

	response.Metadata.CurrentPage = metadata.CurrentPage
	response.Metadata.PageSize = metadata.PageSize
	response.Metadata.FirstPage = metadata.FirstPage
	response.Metadata.LastPage = metadata.LastPage

	if err := platformhandler.WriteJSON(w, http.StatusOK, response, nil); err != nil {
		platformErr.ServerErrorResponse(ctx, w, h.logger, err)
	}
}

func ValidateFilters(v *validator.Validator, f entity.Filters) {
	// Check that the page and page_size parameters contain sensible values. v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10000000, "page", []string{"must be a maximum of 10 million"})
	v.Check(f.PageSize > 0, "pageSize", []string{"must be greater than zero"})
	v.Check(f.Page > 0, "page", []string{"must be greater than zero"})
	v.Check(f.PageSize <= 100, "pageSize", []string{"must be a maximum of 100"})
	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", []string{"invalid sort value"})
}

func mapMasterCardData(a entity.Authorization) *MasterCardDataResponse {
	return &MasterCardDataResponse{
		DE61: DE61{
			TerminalAttendance:                       a.MastercardSchemeData.Request.PointOfServiceData.TerminalAttendance,
			TerminalLocation:                         a.MastercardSchemeData.Request.PointOfServiceData.TerminalLocation,
			CardHolderPresence:                       a.MastercardSchemeData.Request.PointOfServiceData.CardHolderPresence,
			CardPresence:                             a.MastercardSchemeData.Request.PointOfServiceData.CardPresence,
			CardCaptureCapabilities:                  a.MastercardSchemeData.Request.PointOfServiceData.CardCaptureCapabilities,
			TransactionStatus:                        a.MastercardSchemeData.Request.PointOfServiceData.TransactionStatus,
			TransactionSecurity:                      a.MastercardSchemeData.Request.PointOfServiceData.TransactionSecurity,
			CardHolderActivatedTerminalLevel:         a.MastercardSchemeData.Request.PointOfServiceData.CardHolderActivatedTerminalLevel,
			CardDataTerminalInputCapabilityIndicator: a.MastercardSchemeData.Request.PointOfServiceData.CardDataTerminalInputCapabilityIndicator,
			AuthorizationLifeCycle:                   a.MastercardSchemeData.Request.PointOfServiceData.AuthorizationLifeCycle,
			CountryCode:                              a.MastercardSchemeData.Request.PointOfServiceData.CountryCode,
			PostalCode:                               a.MastercardSchemeData.Request.PointOfServiceData.PostalCode,
		},
		DE22: DE22{
			SF1: a.CardSchemeData.Request.POSEntryMode.PanEntryMode.String(),
		},
		DE3: a.CardSchemeData.Request.ProcessingCode.TransactionTypeCode + a.CardSchemeData.Request.ProcessingCode.FromAccountTypeCode + a.CardSchemeData.Request.ProcessingCode.TransactionTypeCode,
	}
}

func mapGetAuthorizationResponse(a entity.Authorization) authorizationResponse {
	response := mapAuthorizationResponse(a)
	if a.Card.Info.Scheme == string(entity.Mastercard) {
		response.MasterCardData = mapMasterCardData(a)
	}
	return response
}
