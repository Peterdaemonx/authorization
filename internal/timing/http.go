package timing

import (
	"gitlab.cmpayments.local/libraries-go/http/jsonresult"
	"net/http"
)

func MetricsHandler(t *Tracker) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Copy the rings; prevent race conditions
		rings := []ring{}
		for _, rp := range t.rings {
			//nolint:govet
			rings = append(rings, *rp)
		}

		res := map[string]map[string]*Timing{}
		//nolint:govet
		for _, r := range rings {
			total, tags := r.Summarize("scheme.mastercard.authorize", "scheme.visa.authorize")
			res[r.ShortLabel()] = tags
			res[r.ShortLabel()]["total"] = total
		}

		//log.Print("%#v", rings)
		jsonresult.Ok(writer, res)
	}
}
