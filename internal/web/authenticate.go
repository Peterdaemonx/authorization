package web

import (
	"context"
	"errors"
	"net/http"

	"gitlab.cmpayments.local/creditcard/platform"
	"gitlab.cmpayments.local/libraries-go/logging"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"

	"gitlab.cmpayments.local/libraries-go/http/jsonresult"
)

type PermissionStore interface {
	GetPermissionsForAPIKey(ctx context.Context, apiKey string) ([]processing.Permission, error)
	//GetPspForAPIKey(ctx context.Context, apiKey string) (processing.PSP, error)
	GetPermissionsForAccountGuid(ctx context.Context, accountGuid string) ([]processing.Permission, error)
}

var ErrUnauthenticated = errors.New("unauthenticated")

type PlatformClient interface {
	GetAccount(ctx context.Context, authCookies []*http.Cookie) (PlatformAccount, error)
}

type PlatformAccount interface {
	IsEmployee() bool
	PersonGuid() string
}

// RequestAuthenticator provides authentication for incoming HTTP requests by wrapping http.HanderFuncs.
// Requests can be authenticated by sending in an API key via the username portion of HTTP Basic Auth,
// and by sending a CMAuthentication cookie.
//
// An API key will be checked against the internal database. CMAuthentication cookies will be used to query the
// CM Identify API to ensure it is valid and belongs to an employee; from there the PersonGuid is used to query the
// internal database for the exact permissions.
type RequestAuthenticator struct {
	Psps PspStore
	Ps   PermissionStore
	Pc   PlatformClient
	Log  platform.Logger
}

func NewRequestAuthenticator(
	psps PspStore,
	ps PermissionStore,
	pc PlatformClient,
	log platform.Logger,
) RequestAuthenticator {
	return RequestAuthenticator{
		Psps: psps,
		Pc:   pc,
		Ps:   ps,
		Log:  log,
	}
}

// WithPermission wraps a http.HanderFunc to check that the credentials in the request
// provide the required permission.
func (ra RequestAuthenticator) WithPermission(perm string, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var permissions []processing.Permission
		ctx := req.Context()
		apiKey, ok := GetAPiKey(req)
		if ok {
			if len([]byte(apiKey)) < 32 {
				ra.Log.Error(ctx, "invalid API key")
				jsonresult.Forbidden(rw)
				return
			}

			perms, err := ra.Ps.GetPermissionsForAPIKey(ctx, apiKey)
			if err != nil {
				ra.Log.Error(logging.ContextWithError(ctx, err), "failed to get permissions for API key")
				jsonresult.InternalServerError(rw, "")
				return
			}
			permissions = append(permissions, perms...)

			psp, err := ra.Psps.GetPspByAPIKey(ctx, apiKey)
			if err != nil {
				ra.Log.Error(logging.ContextWithError(ctx, err), "failed to get psp by API key")
				jsonresult.InternalServerError(rw, "")
				return
			}
			ctx = processing.ContextWithPsp(ctx, psp)
		}

		accountGuid, ok, err := ra.getAccountGuid(req)
		if err != nil {
			ra.Log.Error(logging.ContextWithError(ctx, err), "failed to get account-guid for request")
			jsonresult.InternalServerError(rw, "")
			return
		}
		if ok {
			perms, err := ra.Ps.GetPermissionsForAccountGuid(ctx, accountGuid)
			if err != nil {
				ra.Log.Error(logging.ContextWithError(ctx, err), "failed to get permissions for account-guid")
				// TODO Use logging component
				jsonresult.InternalServerError(rw, "")
				return
			}
			permissions = append(permissions, perms...)
		}

		if !isPermGranted(perm, permissions) {
			jsonresult.Forbidden(rw)
			return
		}

		handler.ServeHTTP(rw, req.WithContext(ctx))
	}
}

func isPermGranted(perm string, perms []processing.Permission) bool {
	for _, check := range perms {
		if perm == check.Code {
			return true
		}
	}
	return false
}

func GetAPiKey(r *http.Request) (string, bool) {
	username, _, ok := r.BasicAuth()
	return username, ok
}

func (ra RequestAuthenticator) getAccountGuid(r *http.Request) (string, bool, error) {
	authCookies, err := GetCMAuthCookies(r)
	if err != nil {
		// http.Authorization.Cookie() only ever returns http.ErrCookieNotFound,
		// so we can ignore the error and just say no account guid was found.
		return "", false, nil
	}

	account, err := ra.Pc.GetAccount(r.Context(), authCookies)
	if err != nil && err != ErrUnauthenticated {
		return "", false, err
	}

	// Only CM Employees have access
	if !account.IsEmployee() {
		return "", false, nil
	}

	return account.PersonGuid(), true, nil
}

func GetCMAuthCookies(r *http.Request) ([]*http.Cookie, error) {
	var authCookies []*http.Cookie

	for _, cookie := range r.Cookies() {
		if len(cookie.Name) >= 6 && cookie.Name[0:6] == "CMAuth" {
			authCookies = append(authCookies, cookie)
		}
	}

	if len(authCookies) == 0 {
		return nil, http.ErrNoCookie
	}

	return authCookies, nil
}
