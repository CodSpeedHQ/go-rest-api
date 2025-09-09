package benchmarks

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkAuth_Login_Success benchmarks successful login
func BenchmarkAuth_Login_Success(b *testing.B) {
	setup := NewBenchmarkSetup()
	loginBody := `{"username":"demo","password":"pass"}`

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("POST", "/v1/login", loginBody, nil)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", rec.Code)
		}
	}
}

// BenchmarkAuth_Login_InvalidCredentials benchmarks login with invalid credentials
func BenchmarkAuth_Login_InvalidCredentials(b *testing.B) {
	setup := NewBenchmarkSetup()
	loginBody := `{"username":"invalid","password":"wrong"}`

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("POST", "/v1/login", loginBody, nil)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			b.Fatalf("expected status 401, got %d", rec.Code)
		}
	}
}