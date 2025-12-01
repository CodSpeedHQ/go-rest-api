package pagination

import (
	"net/http"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(1, 100, 1000)
	}
}

func BenchmarkNewFromRequest(b *testing.B) {
	req, _ := http.NewRequest("GET", "/test?page=2&per_page=50", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewFromRequest(req, 1000)
	}
}

func BenchmarkPages_Offset(b *testing.B) {
	pages := New(5, 100, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pages.Offset()
	}
}

func BenchmarkPages_Limit(b *testing.B) {
	pages := New(5, 100, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pages.Limit()
	}
}

func BenchmarkPages_BuildLinks(b *testing.B) {
	pages := New(5, 100, 1000)
	baseURL := "http://example.com/api/items"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pages.BuildLinks(baseURL, 100)
	}
}

func BenchmarkPages_BuildLinkHeader(b *testing.B) {
	pages := New(5, 100, 1000)
	baseURL := "http://example.com/api/items"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pages.BuildLinkHeader(baseURL, 100)
	}
}

func BenchmarkParseInt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseInt("42", 10)
	}
}

func BenchmarkParseIntDefault(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseInt("", 10)
	}
}

func BenchmarkParseIntInvalid(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseInt("invalid", 10)
	}
}
