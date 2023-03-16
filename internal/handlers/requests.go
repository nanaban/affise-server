package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/sync/errgroup"

	"affise-server/internal/config"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrInvalidCountURLs = errors.New("invalid count of URLs")
)

// RequestsHandler represents requests handler.
type RequestsHandler struct {
	config *config.Client
	client *http.Client
}

// NewRequestsHandler creates new instance of requests handler.
func NewRequestsHandler(conf *config.Client) *RequestsHandler {
	return &RequestsHandler{
		config: conf,
		client: &http.Client{
			Timeout: conf.RequestTimeout,
		},
	}
}

// RequestsRequest represents request for requests handler.
type RequestsRequest []string

// validateRequest validates request.
func (h *RequestsHandler) validateRequest(r RequestsRequest) error {
	if len(r) == 0 || len(r) > h.config.RequestMaxURLs {
		return ErrInvalidCountURLs
	}

	return nil
}

// doGET makes GET request to the url and returns response body.
func (h *RequestsHandler) doGET(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return b, nil
}

// doRequests makes GET requests to the urls and returns responses.
func (h *RequestsHandler) doRequests(ctx context.Context, req RequestsRequest) ([][]byte, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(h.config.RequestMaxConcurrent)

	results := make([][]byte, len(req))
	for i, url := range req {
		i, url := i, url
		g.Go(func() error {
			resp, err := h.doGET(ctx, url)
			if err != nil {
				return fmt.Errorf("url %s: %w", url, err)
			}

			results[i] = resp

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

// ServeHTTP implements http.Handler interface.
func (h *RequestsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := ErrMethodNotAllowed
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	var req RequestsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf("decode request: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validateRequest(req); err != nil {
		err = fmt.Errorf("validate request: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := h.doRequests(r.Context(), req)
	if err != nil {
		err = fmt.Errorf("do requests: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		err = fmt.Errorf("encode response: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
