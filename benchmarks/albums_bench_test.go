package benchmarks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkAlbums_GET_Single benchmarks getting a specific album by ID
func BenchmarkAlbums_GET_Single(b *testing.B) {
	setup := NewBenchmarkSetup()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("GET", "/v1/albums/album-1", "", nil)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", rec.Code)
		}
	}
}

// BenchmarkAlbums_GET_List_Small benchmarks listing albums with small page size
func BenchmarkAlbums_GET_List_Small(b *testing.B) {
	setup := NewBenchmarkSetup()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("GET", "/v1/albums?page=1&per_page=10", "", nil)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", rec.Code)
		}
	}
}

// BenchmarkAlbums_GET_List_Large benchmarks listing albums with large page size
func BenchmarkAlbums_GET_List_Large(b *testing.B) {
	setup := NewBenchmarkSetup()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("GET", "/v1/albums?page=1&per_page=50", "", nil)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", rec.Code)
		}
	}
}

// BenchmarkAlbums_GET_List_Paginated benchmarks paginated album listing
func BenchmarkAlbums_GET_List_Paginated(b *testing.B) {
	setup := NewBenchmarkSetup()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	var page int
	for b.Loop() {
		// Simulate pagination through different pages
		page = (page % 10) + 1
		req := makeRequest("GET", fmt.Sprintf("/v1/albums?page=%d&per_page=10", page), "", nil)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", rec.Code)
		}
	}
}

// BenchmarkAlbums_POST_Create benchmarks creating a new album
func BenchmarkAlbums_POST_Create(b *testing.B) {
	setup := NewBenchmarkSetup()
	createBody := `{"name":"Benchmark Album"}`
	headers := makeAuthHeaders()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("POST", "/v1/albums", createBody, headers)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			b.Fatalf("expected status 201, got %d", rec.Code)
		}
	}
}

// BenchmarkAlbums_PUT_Update benchmarks updating an existing album
func BenchmarkAlbums_PUT_Update(b *testing.B) {
	setup := NewBenchmarkSetup()
	updateBody := `{"name":"Updated Benchmark Album"}`
	headers := makeAuthHeaders()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		req := makeRequest("PUT", "/v1/albums/album-1", updateBody, headers)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", rec.Code)
		}
	}
}

// BenchmarkAlbums_DELETE_Remove benchmarks deleting an album
func BenchmarkAlbums_DELETE_Remove(b *testing.B) {
	setup := NewBenchmarkSetup()
	headers := makeAuthHeaders()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	var counter int
	for b.Loop() {
		// Create a unique album ID for each iteration to avoid conflicts
		albumID := fmt.Sprintf("album-%d", (counter%10)+1)
		counter++
		req := makeRequest("DELETE", "/v1/albums/"+albumID, "", headers)
		rec := httptest.NewRecorder()

		setup.Router.ServeHTTP(rec, req)

		// Accept both 200 (success) and 404 (already deleted) as valid responses
		if rec.Code != http.StatusOK && rec.Code != http.StatusNotFound {
			b.Fatalf("expected status 200 or 404, got %d", rec.Code)
		}
	}
}

// BenchmarkAlbums_CRUD_Mixed benchmarks mixed CRUD operations
func BenchmarkAlbums_CRUD_Mixed(b *testing.B) {
	setup := NewBenchmarkSetup()
	headers := makeAuthHeaders()
	createBody := `{"name":"Mixed CRUD Album"}`
	updateBody := `{"name":"Updated Mixed CRUD Album"}`

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	var counter int
	for b.Loop() {
		var req *http.Request
		var expectedStatus int

		// Cycle through different operations based on iteration
		switch counter % 4 {
		case 0: // GET list
			req = makeRequest("GET", "/v1/albums?page=1&per_page=10", "", nil)
			expectedStatus = http.StatusOK
		case 1: // GET single
			req = makeRequest("GET", "/v1/albums/album-1", "", nil)
			expectedStatus = http.StatusOK
		case 2: // POST create
			req = makeRequest("POST", "/v1/albums", createBody, headers)
			expectedStatus = http.StatusCreated
		case 3: // PUT update
			req = makeRequest("PUT", "/v1/albums/album-1", updateBody, headers)
			expectedStatus = http.StatusOK
		}

		rec := httptest.NewRecorder()
		setup.Router.ServeHTTP(rec, req)

		if rec.Code != expectedStatus {
			b.Fatalf("expected status %d, got %d", expectedStatus, rec.Code)
		}
		counter++
	}
}

// BenchmarkAlbums_Concurrent_Read benchmarks concurrent read operations
func BenchmarkAlbums_Concurrent_Read(b *testing.B) {
	setup := NewBenchmarkSetup()

	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("benchmark panicked: %v", r)
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := makeRequest("GET", "/v1/albums?page=1&per_page=20", "", nil)
			rec := httptest.NewRecorder()

			setup.Router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				b.Fatalf("expected status 200, got %d", rec.Code)
			}
		}
	})
}