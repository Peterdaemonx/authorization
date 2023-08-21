package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"gitlab.cmpayments.local/creditcard/authorization/internal/app"
	"gitlab.cmpayments.local/libraries-go/logging"
)

// nolint:unused
func newTestApplication(t *testing.T) *application {
	// We need conf and logger because they are used in middlewares
	c := app.Config{
		Development: struct {
			MockCmPlatform       bool   `yaml:"mock_cm_platform"`
			MockPermissionStore  bool   `yaml:"mock_permission_store"`
			SpannerEmulatorAddr  string `yaml:"spanner_emulator_addr"`
			HumanReadableLogging bool   `yaml:"human_readable_logging"`
			MockData             bool   `yaml:"mock_data"`
			MockTokenization     bool   `yaml:"mock_tokenization"`
			MockCardInfo         bool   `yaml:"mock_card_info"`
			MockPublisher        bool   `yaml:"mock_publisher"`
			MockPubSub           bool   `yaml:"mock_pub_sub"`
			MockNonceStore       bool   `yaml:"mock_nonce_store"`
			MockStorageBucket    bool   `yaml:"mock_storage_bucket"`
		}{
			MockCmPlatform:       true,
			MockPermissionStore:  true,
			SpannerEmulatorAddr:  "",
			HumanReadableLogging: false,
			MockData:             true,
			MockTokenization:     true,
			MockCardInfo:         true,
			MockPublisher:        true,
			MockPubSub:           true,
			MockNonceStore:       true,
			MockStorageBucket:    true,
		},
		Cors: struct {
			AllowedOrigins string `yaml:"allowed_origins"`
		}{
			AllowedOrigins: "",
		},
	}

	return &application{
		ctx:    context.Background(),
		conf:   c,
		logger: logging.Logger{},
	}
}

// nolint:unused
type testServer struct {
	*httptest.Server
}

// nolint:unused
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	// store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// disable redirects after 3xx response codes
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// nolint:unused
func (ts *testServer) get(t *testing.T, urlPath, apiKey string) (int, http.Header, []byte) {
	req, err := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	if apiKey != "" {
		req.SetBasicAuth(apiKey, "")
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)

	}

	return rs.StatusCode, rs.Header, body
}

// nolint:unused
func (ts *testServer) post(t *testing.T, urlPath, apiKey, nonce string, body []byte) (int, http.Header, []byte) {
	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	if apiKey != "" {
		req.SetBasicAuth(apiKey, "")
	}

	if nonce != "" {
		req.Header.Set("nonce", nonce)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	respBody, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)

	}

	return rs.StatusCode, rs.Header, respBody
}
