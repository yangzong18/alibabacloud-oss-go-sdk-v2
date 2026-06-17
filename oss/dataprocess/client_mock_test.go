package dataprocess

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"
)

func sortQuery(r *http.Request) string {
	u := r.URL
	var buf strings.Builder
	keys := make([]string, 0, len(u.Query()))
	for k := range u.Query() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := u.Query()[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			if len(v) > 0 {
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	u.RawQuery = buf.String()
	return u.String()
}

func testSetupMockServer(t *testing.T, statusCode int, headers map[string]string, body []byte,
	chkfunc func(t *testing.T, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check request
		chkfunc(t, r)

		// header s
		for k, v := range headers {
			w.Header().Set(k, v)
		}

		// status code
		w.WriteHeader(statusCode)

		// body
		w.Write(body)
	}))
}
