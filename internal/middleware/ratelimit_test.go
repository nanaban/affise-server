package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval time.Duration
		limit    int
	}{
		{
			name:     "500ms-5",
			interval: 500 * time.Millisecond,
			limit:    5,
		},
		{
			name:     "1s-3",
			interval: time.Second,
			limit:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			middleware := NewRateLimit(tt.interval, tt.limit)
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = fmt.Fprintf(w, "ok")
			})

			server := httptest.NewServer(middleware.Handle(testHandler))
			defer server.Close()

			client := server.Client()

			for i := 0; i < tt.limit; i++ {
				resp, err := client.Get(server.URL)
				if err != nil {
					t.Fatal(err)
				}
				if resp.StatusCode != http.StatusOK {
					t.Errorf("invalid status code: %d", resp.StatusCode)
				}
			}

			// rate limit exceeded
			resp, err := client.Get(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("invalid status code: %d", resp.StatusCode)
			}
		})
	}
}
