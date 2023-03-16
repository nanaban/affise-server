package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"affise-server/internal/config"
)

type testENV struct {
	config *config.Config
	server *Server
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

var env *testENV

func setup(tb testing.TB) {
	tb.Helper()

	conf := config.NewDefault()
	conf.Server.Addr = "localhost:8889"

	ctx, cancel := context.WithCancel(context.Background())
	s := NewServer(conf)

	env = &testENV{
		config: conf,
		server: s,
		ctx:    ctx,
		cancel: cancel,
	}

	env.wg.Add(1)
	go func() {
		defer env.wg.Done()

		if err := s.Run(ctx); err != nil {
			tb.Error(err)
		}
	}()
}

func teardown(tb testing.TB) {
	tb.Helper()

	env.cancel()
	env.wg.Wait()
}

func request(ctx context.Context, method, addr, path string, body any) (*http.Response, []byte, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, fmt.Errorf("couldn't encode body: %w", err)
	}

	endpoint := fmt.Sprintf("http://%s%s", addr, path)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, buf)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't do request: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't read response body: %w", err)
	}

	return resp, b, nil
}

func makeRequest(t *testing.T, method, path string, body any) (*http.Response, []byte) {
	t.Helper()

	resp, b, err := request(env.ctx, method, env.server.Addr(), path, body)
	if err != nil {
		t.Fatal(err)
	}

	return resp, b
}

func repeatToSlice(s string, n int) []string {
	var result []string

	for i := 0; i < n; i++ {
		result = append(result, s)
	}

	return result
}

func TestServer(t *testing.T) {
	setup(t)
	t.Cleanup(func() {
		teardown(t)
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		t.Parallel()

		list := []string{"https://www.google.com"}

		resp, b := makeRequest(t, http.MethodGet, EndPointRequests, list)
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("invalid status code: %d, body: %s", resp.StatusCode, b)
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			list []string
		}{
			{list: []string{}},
			{list: repeatToSlice("https://www.google.com", 21)},
		}

		for _, tt := range tests {
			resp, b := makeRequest(t, http.MethodPost, EndPointRequests, tt.list)
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("invalid status code: %d, body: %s", resp.StatusCode, b)
			}
		}
	})

	t.Run("InvalidBody", func(t *testing.T) {
		t.Parallel()

		list := []int{1, 2, 3}

		resp, b := makeRequest(t, http.MethodPost, EndPointRequests, list)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("invalid status code: %d, body: %s", resp.StatusCode, b)
		}
	})

	t.Run("InvalidURL", func(t *testing.T) {
		t.Parallel()

		list := []string{"https://foo"}

		resp, b := makeRequest(t, http.MethodPost, EndPointRequests, list)
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("invalid status code: %d, body: %s", resp.StatusCode, b)
		}
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		list := repeatToSlice("https://www.google.com", 5)

		resp, b := makeRequest(t, http.MethodPost, EndPointRequests, list)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code: %d, body: %s", resp.StatusCode, b)
		}

		var res [][]byte
		if err := json.Unmarshal(b, &res); err != nil {
			t.Errorf("couldn't unmarshal body: %v", err)
		}
		if len(res) != len(list) {
			t.Errorf("invalid response length: %d", len(res))
		}
	})
}
