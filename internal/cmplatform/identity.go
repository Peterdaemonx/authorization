package cmplatform

import (
	"context"
	"encoding/json"
	"net/http"

	httplog "gitlab.cmpayments.local/creditcard/platform/http/logging"

	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
)

// Define errors
var (
	ErrUnauthenticated = web.ErrUnauthenticated
)

// NewIdentityClient creates new identity client
func NewIdentityClient(domain string, client http.Client) IdentityClient {
	client.Transport = wrapWithJSONAcceptHeaderTransport(client.Transport)
	client.Transport = wrapWithJSONContentTypeHeaderTransport(client.Transport)
	client.Transport = httplog.NewTraceIDRoundTripper(client.Transport)

	return IdentityClient{
		http:           client,
		getIdentityURL: "https://login." + domain + "/v1.0/identity",
		getAccountURL:  "https://api." + domain + "/accounts/v2.0/accounts/",
	}
}

// IdentityClient is used to login and access accounts
type (
	IdentityClient struct {
		http           http.Client
		getIdentityURL string
		getAccountURL  string
	}

	identityResponse struct {
		Code      int
		Message   string
		SignInURL string
		UserData  struct {
			PersonGUID string
		}
		Claims struct {
			IsEmployee           bool   `json:"IsEmployee"`
			AuthenticationMethod string `json:"authenticationmethod"`
		}
	}
)

func (c IdentityClient) GetAccount(ctx context.Context, authCookies []*http.Cookie) (web.PlatformAccount, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getIdentityURL, nil)
	if err != nil {
		return identityResponse{}, err
	}

	for _, cookie := range authCookies {
		request.AddCookie(cookie)
	}

	resp, err := c.http.Do(request)
	if err != nil {
		return identityResponse{}, err
	}
	response := identityResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return identityResponse{}, err
	}
	if response.Code == 100 {
		return identityResponse{}, ErrUnauthenticated
	}
	return response, nil
}

func (ir identityResponse) PersonGuid() string {
	return ir.UserData.PersonGUID
}

func (ir identityResponse) IsEmployee() bool {
	return ir.Claims.IsEmployee || ir.Claims.AuthenticationMethod == "ActiveDirectory"
}
