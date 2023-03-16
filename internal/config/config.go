package config

import "time"

const (
	DefaultServerAddr                 = ":8080"
	DefaultServerRateInterval         = 1 * time.Second
	DefaultServerRateLimit            = 100
	DefaultClientRequestTimeout       = 1 * time.Second
	DefaultClientRequestMaxConcurrent = 4
	DefaultClientRequestMaxURLs       = 20
)

// Server represents server configuration.
type Server struct {
	Addr         string        `env:"SERVER_ADDR"`
	RateLimit    int           `env:"SERVER_RATE_LIMIT"`
	RateInterval time.Duration `env:"SERVER_RATE_INTERVAL"`
}

// Client represents client configuration.
type Client struct {
	RequestTimeout       time.Duration `env:"CLIENT_REQUEST_TIMEOUT"`
	RequestMaxConcurrent int           `env:"CLIENT_REQUEST_MAX_CONCURRENT"`
	RequestMaxURLs       int           `env:"CLIENT_REQUEST_MAX_URLS"`
}

// Config represents configuration.
type Config struct {
	Server Server
	Client Client
}

// NewDefault creates new instance of configuration with default values.
func NewDefault() *Config {
	c := &Config{}
	c.SetDefaults()
	return c
}

// SetDefaults sets default values.
func (c *Config) SetDefaults() {
	if c.Server.Addr == "" {
		c.Server.Addr = DefaultServerAddr
	}
	if c.Server.RateLimit == 0 {
		c.Server.RateLimit = DefaultServerRateLimit
	}
	if c.Server.RateInterval == 0 {
		c.Server.RateInterval = DefaultServerRateInterval
	}
	if c.Client.RequestTimeout == 0 {
		c.Client.RequestTimeout = DefaultClientRequestTimeout
	}
	if c.Client.RequestMaxConcurrent == 0 {
		c.Client.RequestMaxConcurrent = DefaultClientRequestMaxConcurrent
	}
	if c.Client.RequestMaxURLs == 0 {
		c.Client.RequestMaxURLs = DefaultClientRequestMaxURLs
	}
}
