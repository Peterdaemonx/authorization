//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestGetAuthorizations(t *testing.T) {
	app := newTestApplication(t)

	r, err := app.routes()
	if err != nil {
		t.Fatal(err)
	}

	ts := newTestServer(t, r)
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{
			name:     "getAuthorizationsMetadata",
			urlPath:  "/v1/authorizations?page=1",
			wantCode: http.StatusOK,
			//wantBody: []byte(getAuthorizationsMetadata),
		},
		{
			name:     "getAuthorizations",
			urlPath:  "/v1/authorizations?ref=ABC123&amount=10000&date=2022-02-10&pan=1411&status=approved&exemption=mit&responseCode=1&responseId=1&page=1&pageSize=15",
			wantCode: http.StatusOK,
			//wantBody: []byte(getAuthorizations),
		},
		{
			name:     "getAuthorizationWithExemption",
			urlPath:  "/v1/authorizations?ref=ABC123&amount=10000&date=2022-02-10&pan=1411&status=approved&exemption=mit&responseCode=1&responseId=1&page=1&pageSize=15",
			wantCode: http.StatusOK,
			//wantBody: []byte(getAuthorizationWithExemption),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath, "6247c10c-84a0-4fa1-b330-77eea1e944d3")

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("\nwant body to contain %q\ngot %q", tt.wantBody, body)
			}
		})
	}
}

func TestPostCaptures(t *testing.T) {
	app := newTestApplication(t)

	r, err := app.routes()

	if err != nil {
		t.Fatal(err)
	}

	ts := newTestServer(t, r)
	defer ts.Close()

	tests := []struct {
		name           string
		urlPath        string
		reqBody        []byte
		wantCode       int
		wantBody       []byte
		wantProperties map[string]interface{}
	}{
		{name: "captureAuthorization",
			urlPath:  "/v1/authorizations/750b17dd-b89b-4991-b3d7-cca78ca7a654/captures",
			reqBody:  []byte("{\"amount\": 100,\"isFinal\": true, \"currency\": \"EUR\"}"),
			wantCode: http.StatusCreated,
			wantBody: []byte("{\"id\":\"6b73e47e-e8d5-4949-8b2e-31fa49e814ab\",\"authorizationId\":\"750b17dd-b89b-4991-b3d7-cca78ca7a654\",\"url\":\"/v1/authorizations/750b17dd-b89b-4991-b3d7-cca78ca7a654/capture/6b73e47e-e8d5-4949-8b2e-31fa49e814ab\",\"amount\":100,\"currency\":\"EUR\",\"isFinal\":true}"),
			wantProperties: map[string]interface{}{
				"authorizationId": "750b17dd-b89b-4991-b3d7-cca78ca7a654",
				"amount":          float64(100),
				"currency":        "EUR",
				"isFinal":         true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.post(t, tt.urlPath, "6247c10c-84a0-4fa1-b330-77eea1e944d3", "gw-7659fad6288c", tt.reqBody)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			props := make(map[string]interface{})
			err := json.Unmarshal(body, &props)
			if err != nil {
				t.Errorf("cannot parse response: %s", err.Error())
			}

			for k, v := range tt.wantProperties {
				if !reflect.DeepEqual(props[k], v) {
					t.Errorf("\nwant\t%#v\ngot\t\t%#v", v, props[k])
				}
			}
		})
	}
}
