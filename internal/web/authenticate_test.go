package web_test

//go:generate mockgen -package=web_test -source=./authenticate.go -destination=./authenticate_mock_test.go
//go:generate mockgen -package=web_test -source=./nonce.go -destination=./nonce_mock_test.go
//go:generate mockgen -package=web_test -destination=./authenticate_mock_http_test.go net/http Handler,ResponseWriter

import (
	"errors"
	"net/http"
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"

	"github.com/google/uuid"
	storage_mocks "gitlab.cmpayments.local/creditcard/authorization/internal/processing/mocks"
	"gitlab.cmpayments.local/libraries-go/logging"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"

	"github.com/golang/mock/gomock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
)

func TestRequestAuthenticator_WithPermission_APIKey_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	pspStore := storage_mocks.NewMockPspStore(ctrl)
	pspStore.EXPECT().GetPspByAPIKey(gomock.Any(), "f6be7e92-657f-425b-8ce4-fcdfe84c313c").Return(entity.PSP{ID: uuid.New(), Name: "test_psp", Prefix: "tps"}, nil)

	permStore := NewMockPermissionStore(ctrl)
	permStore.EXPECT().GetPermissionsForAPIKey(gomock.Any(), "f6be7e92-657f-425b-8ce4-fcdfe84c313c").Return([]processing.Permission{{Code: "create_authorization"}}, nil)

	platformClient := NewMockPlatformClient(ctrl)

	handler := NewMockHandler(ctrl)
	handler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())

	reqAuth := web.RequestAuthenticator{
		Psps: pspStore,
		Ps:   permStore,
		Pc:   platformClient,
	}

	wrappedHandler := reqAuth.WithPermission("create_authorization", handler.ServeHTTP)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("f6be7e92-657f-425b-8ce4-fcdfe84c313c", "")
	rw := NewMockResponseWriter(ctrl)

	wrappedHandler(rw, req)
}

func TestRequestAuthenticator_WithPermission_APIKey_PermissionStoreError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permStore := NewMockPermissionStore(ctrl)
	permStore.EXPECT().GetPermissionsForAPIKey(gomock.Any(), "f6be7e92-657f-425b-8ce4-fcdfe84c313c").Return(nil, errors.New("test"))

	platformClient := NewMockPlatformClient(ctrl)

	handler := NewMockHandler(ctrl)

	reqAuth := web.RequestAuthenticator{
		Ps:  permStore,
		Pc:  platformClient,
		Log: logging.Logger{},
	}

	wrappedHandler := reqAuth.WithPermission("create_authorization", handler.ServeHTTP)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("f6be7e92-657f-425b-8ce4-fcdfe84c313c", "")
	rw := NewMockResponseWriter(ctrl)
	rw.EXPECT().WriteHeader(500)
	rw.EXPECT().Header().Return(http.Header{})
	rw.EXPECT().Write(gomock.Any())

	wrappedHandler(rw, req)
}

func TestRequestAuthenticator_WithPermission_Cookie_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pspStore := storage_mocks.NewMockPspStore(ctrl)
	pspStore.EXPECT().GetPspByAPIKey(gomock.Any(), "f6be7e92-657f-425b-8ce4-fcdfe84c313c").Return(entity.PSP{ID: uuid.New(), Name: "test_psp", Prefix: "tps"}, nil)

	permStore := NewMockPermissionStore(ctrl)
	permStore.EXPECT().GetPermissionsForAccountGuid(gomock.Any(), "1234").Return([]processing.Permission{{Code: "create_authorization"}}, nil)
	permStore.EXPECT().GetPermissionsForAPIKey(gomock.Any(), "f6be7e92-657f-425b-8ce4-fcdfe84c313c").Return([]processing.Permission{{Code: "create_authorization"}}, nil)
	account := NewMockPlatformAccount(ctrl)
	account.EXPECT().PersonGuid().Return("1234")
	account.EXPECT().IsEmployee().Return(true)

	cookies := []*http.Cookie{
		{Name: "CMAuthentication", Value: "0987654321"},
	}

	platformClient := NewMockPlatformClient(ctrl)
	platformClient.EXPECT().GetAccount(gomock.Any(), cookies).Return(account, nil)

	handler := NewMockHandler(ctrl)
	handler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())

	reqAuth := web.RequestAuthenticator{
		Psps: pspStore,
		Ps:   permStore,
		Pc:   platformClient,
	}

	wrappedHandler := reqAuth.WithPermission("create_authorization", handler.ServeHTTP)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("f6be7e92-657f-425b-8ce4-fcdfe84c313c", "")

	if err != nil {
		t.Fail()
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	rw := NewMockResponseWriter(ctrl)
	wrappedHandler(rw, req)
}

func TestRequestAuthenticator_WithPermission_Cookie_PermissionStoreError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	pspID := uuid.New()
	pspStore := storage_mocks.NewMockPspStore(ctrl)
	pspStore.EXPECT().GetPspByAPIKey(gomock.Any(), gomock.Any()).Return(entity.PSP{ID: pspID, Name: "test_psp", Prefix: "tps"}, nil)

	permStore := NewMockPermissionStore(ctrl)
	permStore.EXPECT().GetPermissionsForAccountGuid(gomock.Any(), gomock.Any()).Return([]processing.Permission{}, nil)
	permStore.EXPECT().GetPermissionsForAPIKey(gomock.Any(), "f6be7e92-657f-425b-8ce4-fcdfe84c313c").Return([]processing.Permission{}, nil)
	account := NewMockPlatformAccount(ctrl)
	account.EXPECT().PersonGuid().Return("f6be7e92-657f-425b-8ce4-fcdfe84c313c")
	account.EXPECT().IsEmployee().Return(true)

	cookies := []*http.Cookie{
		{Name: "CMAuthentication", Value: "0987654321"},
	}

	platformClient := NewMockPlatformClient(ctrl)
	platformClient.EXPECT().GetAccount(gomock.Any(), cookies).Return(account, nil)

	handler := NewMockHandler(ctrl)

	reqAuth := web.RequestAuthenticator{
		Psps: pspStore,
		Ps:   permStore,
		Pc:   platformClient,
	}

	wrappedHandler := reqAuth.WithPermission("create_authorization", handler.ServeHTTP)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	ctx := processing.ContextWithPsp(req.Context(), entity.PSP{ID: pspID})
	req = req.WithContext(ctx)
	req.SetBasicAuth("f6be7e92-657f-425b-8ce4-fcdfe84c313c", "")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rw := NewMockResponseWriter(ctrl)
	rw.EXPECT().WriteHeader(403)

	wrappedHandler(rw, req)
}

func TestRequestAuthenticator_WithPermission_Cookie_PlatformClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pspStore := storage_mocks.NewMockPspStore(ctrl)
	cookies := []*http.Cookie{
		{Name: "CMAuthentication", Value: "0987654321"},
	}
	account := NewMockPlatformAccount(ctrl)
	account.EXPECT().IsEmployee().Return(true)
	account.EXPECT().PersonGuid().Return("1234")
	platformClient := NewMockPlatformClient(ctrl)
	platformClient.EXPECT().GetAccount(gomock.Any(), cookies).Return(account, nil)
	permStore := NewMockPermissionStore(ctrl)
	permStore.EXPECT().GetPermissionsForAccountGuid(gomock.Any(), "1234").Return([]processing.Permission{}, nil)
	handler := NewMockHandler(ctrl)

	reqAuth := web.RequestAuthenticator{
		Psps: pspStore,
		Ps:   permStore,
		Pc:   platformClient,
	}

	wrappedHandler := reqAuth.WithPermission("create_authorization", handler.ServeHTTP)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rw := NewMockResponseWriter(ctrl)

	rw.EXPECT().WriteHeader(403)
	wrappedHandler(rw, req)
}

func TestRequestAuthenticator_WithPermission_Cookie_NotEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pspStore := storage_mocks.NewMockPspStore(ctrl)

	permStore := NewMockPermissionStore(ctrl)

	account := NewMockPlatformAccount(ctrl)
	account.EXPECT().IsEmployee().Return(false)

	cookies := []*http.Cookie{
		{Name: "CMAuthentication", Value: "0987654321"},
	}

	platformClient := NewMockPlatformClient(ctrl)
	platformClient.EXPECT().GetAccount(gomock.Any(), cookies).Return(account, nil)

	handler := NewMockHandler(ctrl)

	reqAuth := web.RequestAuthenticator{
		Psps: pspStore,
		Ps:   permStore,
		Pc:   platformClient,
	}

	wrappedHandler := reqAuth.WithPermission("create_authorization", handler.ServeHTTP)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rw := NewMockResponseWriter(ctrl)
	rw.EXPECT().WriteHeader(403)

	wrappedHandler(rw, req)
}
