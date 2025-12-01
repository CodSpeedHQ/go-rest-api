package log

import (
	"context"
	"net/http"
	"testing"
)

func BenchmarkLogger_With(b *testing.B) {
	logger := New()
	ctx := context.Background()
	args := []interface{}{"key1", "value1", "key2", "value2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logger.With(ctx, args...)
	}
}

func BenchmarkLogger_WithContext(b *testing.B) {
	logger := New()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	req.Header.Set("X-Correlation-ID", "test-correlation-id")
	ctx := WithRequest(context.Background(), req)
	args := []interface{}{"key", "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logger.With(ctx, args...)
	}
}

func BenchmarkLogger_Debug(b *testing.B) {
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("debug message")
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("info message")
	}
}

func BenchmarkLogger_Error(b *testing.B) {
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("error message")
	}
}

func BenchmarkLogger_Debugf(b *testing.B) {
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debugf("debug message: %s %d", "test", i)
	}
}

func BenchmarkLogger_Infof(b *testing.B) {
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("info message: %s %d", "test", i)
	}
}

func BenchmarkLogger_Errorf(b *testing.B) {
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Errorf("error message: %s %d", "test", i)
	}
}

func BenchmarkWithRequest(b *testing.B) {
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	req.Header.Set("X-Correlation-ID", "test-correlation-id")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithRequest(ctx, req)
	}
}
