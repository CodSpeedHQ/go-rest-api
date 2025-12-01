package album

import (
	"context"
	"testing"

	"github.com/qiangxue/go-rest-api/pkg/log"
)

func BenchmarkService_Create(b *testing.B) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)
	ctx := context.Background()
	req := CreateAlbumRequest{Name: "benchmark album"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Create(ctx, req)
	}
}

func BenchmarkService_Get(b *testing.B) {
	logger, _ := log.NewForTest()
	repo := &mockRepository{}
	s := NewService(repo, logger)
	ctx := context.Background()

	// Setup: create an album first
	album, _ := s.Create(ctx, CreateAlbumRequest{Name: "test album"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Get(ctx, album.ID)
	}
}

func BenchmarkService_Update(b *testing.B) {
	logger, _ := log.NewForTest()
	repo := &mockRepository{}
	s := NewService(repo, logger)
	ctx := context.Background()

	// Setup: create an album first
	album, _ := s.Create(ctx, CreateAlbumRequest{Name: "test album"})
	req := UpdateAlbumRequest{Name: "updated album"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Update(ctx, album.ID, req)
	}
}

func BenchmarkService_Delete(b *testing.B) {
	logger, _ := log.NewForTest()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		repo := &mockRepository{}
		s := NewService(repo, logger)
		album, _ := s.Create(ctx, CreateAlbumRequest{Name: "test album"})
		b.StartTimer()

		_, _ = s.Delete(ctx, album.ID)
	}
}

func BenchmarkService_Query(b *testing.B) {
	logger, _ := log.NewForTest()
	repo := &mockRepository{}
	s := NewService(repo, logger)
	ctx := context.Background()

	// Setup: create multiple albums
	for i := 0; i < 10; i++ {
		_, _ = s.Create(ctx, CreateAlbumRequest{Name: "test album"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Query(ctx, 0, 10)
	}
}

func BenchmarkService_Count(b *testing.B) {
	logger, _ := log.NewForTest()
	repo := &mockRepository{}
	s := NewService(repo, logger)
	ctx := context.Background()

	// Setup: create multiple albums
	for i := 0; i < 10; i++ {
		_, _ = s.Create(ctx, CreateAlbumRequest{Name: "test album"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Count(ctx)
	}
}

func BenchmarkCreateAlbumRequest_Validate(b *testing.B) {
	req := CreateAlbumRequest{Name: "test album"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = req.Validate()
	}
}

func BenchmarkUpdateAlbumRequest_Validate(b *testing.B) {
	req := UpdateAlbumRequest{Name: "updated album"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = req.Validate()
	}
}
