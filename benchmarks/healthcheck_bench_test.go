package benchmarks

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkHealthcheck_GET benchmarks the health check endpoint
func BenchmarkHealthcheck_GET(b *testing.B) {
	setup := NewBenchmarkSetup()
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("GET", "/healthcheck", "", nil)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", rec.Code)
		}
	}
}