package cmplatform

import "net/http"

func wrapWithJSONAcceptHeaderTransport(next http.RoundTripper) http.RoundTripper {
	return jsonAcceptHeaderTransport{next: next}
}

type jsonAcceptHeaderTransport struct {
	next http.RoundTripper
}

func (tp jsonAcceptHeaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Header.Get("Accept") == "" {
		r.Header.Set("Accept", "application/json")
	}
	if tp.next == nil {
		return http.DefaultTransport.RoundTrip(r)
	}
	return tp.next.RoundTrip(r)
}
