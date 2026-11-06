package platform

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	h := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}), 2, time.Minute, "/api/leads")
	call := func(method, path, ip string) int {
		req := httptest.NewRequest(method, path, nil)
		req.Header.Set("X-Forwarded-For", ip+", 10.0.0.1")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}
	for i, want := range []int{200, 200, 429, 429} {
		if got := call("POST", "/api/leads", "1.1.1.1"); got != want {
			t.Fatalf("post %d: %d, want %d", i+1, got, want)
		}
	}
	if call("POST", "/api/leads", "2.2.2.2") != 200 {
		t.Error("another IP has its own limit")
	}
	if call("GET", "/api/leads", "1.1.1.1") != 200 || call("POST", "/api/orders", "1.1.1.1") != 200 {
		t.Error("only POSTs to the listed paths count")
	}
}
