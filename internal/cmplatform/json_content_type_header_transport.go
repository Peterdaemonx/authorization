package cmplatform

import "net/http"

func wrapWithJSONContentTypeHeaderTransport(next http.RoundTripper) http.RoundTripper {
	return jsonContentTypeHeaderTransport{next: next}
}

type jsonContentTypeHeaderTransport struct {
	next http.RoundTripper
}

func (tp jsonContentTypeHeaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Body != nil && r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", "application/json")
	}
	if tp.next == nil {
		return http.DefaultTransport.RoundTrip(r)
	}
	return tp.next.RoundTrip(r)
}
