package auth

import (
	"context"
	"testing"

	"github.com/qiangxue/go-rest-api/pkg/log"
)

func BenchmarkService_Login(b *testing.B) {
	logger, _ := log.NewForTest()
	s := NewService("test-signing-key", 3600, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Login(ctx, "demo", "pass")
	}
}

func BenchmarkService_LoginInvalid(b *testing.B) {
	logger, _ := log.NewForTest()
	s := NewService("test-signing-key", 3600, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Login(ctx, "invalid", "invalid")
	}
}

func BenchmarkWithUser(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithUser(ctx, "100", "demo")
	}
}

func BenchmarkCurrentUser(b *testing.B) {
	ctx := WithUser(context.Background(), "100", "demo")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CurrentUser(ctx)
	}
}
