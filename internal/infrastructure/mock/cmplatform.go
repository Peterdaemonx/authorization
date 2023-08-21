package mock

import (
	"context"
	"net/http"
	"strings"

	"gitlab.cmpayments.local/creditcard/authorization/internal/web"
)

type PlatformClient struct{}

func (PlatformClient) GetAccount(ctx context.Context, authCookies []*http.Cookie) (web.PlatformAccount, error) {
	val := strings.SplitN(authCookies[0].Value, "_", 2)
	return platformAcount{
		guid:       val[0],
		isEmployee: len(val) > 1 && val[1] == "employee",
	}, nil
}

type platformAcount struct {
	guid       string
	isEmployee bool
}

func (p platformAcount) IsEmployee() bool {
	return p.isEmployee
}

func (p platformAcount) PersonGuid() string {
	return p.guid
}
