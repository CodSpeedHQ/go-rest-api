package benchmarks

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/qiangxue/go-rest-api/internal/album"
	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/internal/healthcheck"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// BenchmarkSetup provides shared setup utilities for benchmarks
type BenchmarkSetup struct {
	Logger log.Logger
	Router *routing.Router
}

// NewBenchmarkSetup creates a new benchmark setup with all necessary components
func NewBenchmarkSetup() *BenchmarkSetup {
	logger, _ := log.NewForTest()
	router := buildBenchmarkRouter(logger)
	return &BenchmarkSetup{
		Logger: logger,
		Router: router,
	}
}

// buildBenchmarkRouter creates a router similar to the production setup but optimized for benchmarks
func buildBenchmarkRouter(logger log.Logger) *routing.Router {
	router := routing.New()

	// Use minimal middleware for benchmarks to reduce overhead
	router.Use(
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	// Register health check
	healthcheck.RegisterHandlers(router, "benchmark")

	// Create v1 group
	rg := router.Group("/v1")

	// Set up album handlers with mock repository
	albumRepo := &mockAlbumRepository{
		items: createMockAlbums(100), // Pre-populate with 100 albums
	}
	albumService := album.NewService(albumRepo, logger)
	
	// Create a new route group for albums to ensure proper middleware setup
	albumGroup := rg.Group("")
	album.RegisterHandlers(albumGroup, albumService, auth.MockAuthHandler, logger)

	// Set up auth handlers with mock service
	authService := &mockAuthService{}
	auth.RegisterHandlers(rg.Group(""), authService, logger)

	return router
}

// mockAlbumRepository implements album.Repository for benchmarks
type mockAlbumRepository struct {
	items   []entity.Album
	counter int
}

func (r *mockAlbumRepository) Get(ctx context.Context, id string) (entity.Album, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Album{}, errors.NotFound("")
}

func (r *mockAlbumRepository) Count(ctx context.Context) (int, error) {
	return len(r.items), nil
}

func (r *mockAlbumRepository) Query(ctx context.Context, offset, limit int) ([]entity.Album, error) {
	if offset >= len(r.items) {
		return []entity.Album{}, nil
	}
	end := offset + limit
	if end > len(r.items) {
		end = len(r.items)
	}
	return r.items[offset:end], nil
}

func (r *mockAlbumRepository) Create(ctx context.Context, album entity.Album) error {
	// The album ID is already set by the service, just store it
	r.items = append(r.items, album)
	return nil
}

func (r *mockAlbumRepository) Update(ctx context.Context, album entity.Album) error {
	for i, item := range r.items {
		if item.ID == album.ID {
			album.UpdatedAt = time.Now()
			r.items[i] = album
			return nil
		}
	}
	return errors.NotFound("")
}

func (r *mockAlbumRepository) Delete(ctx context.Context, id string) error {
	for i, item := range r.items {
		if item.ID == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.NotFound("")
}

// mockAuthService implements auth.Service for benchmarks
type mockAuthService struct{}

func (s *mockAuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username == "demo" && password == "pass" {
		return "benchmark-token", nil
	}
	return "", errors.Unauthorized("")
}

// createMockAlbums generates test album data for benchmarks
func createMockAlbums(count int) []entity.Album {
	albums := make([]entity.Album, count)
	now := time.Now()
	for i := 0; i < count; i++ {
		albums[i] = entity.Album{
			ID:        fmt.Sprintf("album-%d", i+1),
			Name:      fmt.Sprintf("Benchmark Album %d", i+1),
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return albums
}

// makeRequest is a helper function to create HTTP requests for benchmarks
func makeRequest(method, url string, body string, headers http.Header) *http.Request {
	var bodyReader *strings.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
		req, _ := http.NewRequest(method, url, bodyReader)
		if headers != nil {
			req.Header = headers
		}
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
		return req
	}

	req, _ := http.NewRequest(method, url, nil)
	if headers != nil {
		req.Header = headers
	}
	return req
}

// makeAuthHeaders creates authentication headers for protected endpoints
func makeAuthHeaders() http.Header {
	return auth.MockAuthHeader()
}