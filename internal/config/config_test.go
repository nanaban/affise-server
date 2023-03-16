package config

import (
	"testing"
)

func TestNewDefault(t *testing.T) {
	c := NewDefault()
	if c.Server.Addr != DefaultServerAddr {
		t.Errorf("expected %s, got %s", DefaultServerAddr, c.Server.Addr)
	}
	if c.Server.RateLimit != DefaultServerRateLimit {
		t.Errorf("expected %d, got %d", DefaultServerRateLimit, c.Server.RateLimit)
	}
	if c.Server.RateInterval != DefaultServerRateInterval {
		t.Errorf("expected %s, got %s", DefaultServerRateInterval, c.Server.RateInterval)
	}
	if c.Client.RequestTimeout != DefaultClientRequestTimeout {
		t.Errorf("expected %s, got %s", DefaultClientRequestTimeout, c.Client.RequestTimeout)
	}
	if c.Client.RequestMaxConcurrent != DefaultClientRequestMaxConcurrent {
		t.Errorf("expected %d, got %d", DefaultClientRequestMaxConcurrent, c.Client.RequestMaxConcurrent)
	}
	if c.Client.RequestMaxURLs != DefaultClientRequestMaxURLs {
		t.Errorf("expected %d, got %d", DefaultClientRequestMaxURLs, c.Client.RequestMaxURLs)
	}
}
