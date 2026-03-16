package cache

import (
	"errors"
	"time"
)

// ErrInvalidCapacity is returned by WithCapacity when n is less than or equal to zero.
var ErrInvalidCapacity = errors.New("cap must be greater than zero")

// ErrInvalidTTL is returned by WithTTL when d is less than or equal to zero.
var ErrInvalidTTL = errors.New("ttl must be greater than zero")

// Config holds the configuration parameters for an LRU cache.
// Use Option functions to populate it; do not construct it directly.
type Config struct {
	capacity int
	ttl      time.Duration
}

// Option is a functional option that configures a [Config].
type Option func(*Config) error

// WithCapacity sets the maximum number of entries the cache can hold.
// Returns [ErrInvalidCapacity] if n is less than or equal to zero.
func WithCapacity(n int) Option {
	return func(c *Config) error {
		if n <= 0 {
			return ErrInvalidCapacity
		}
		c.capacity = n
		return nil
	}
}

// WithTTL sets the time-to-live for each cache entry.
// Entries are considered expired after d elapses since their last write and
// will not be returned by Get or Exists. Expired entries are evicted lazily on access.
// Returns [ErrInvalidTTL] if d is less than or equal to zero.
func WithTTL(d time.Duration) Option {
	return func(c *Config) error {
		if d <= 0 {
			return ErrInvalidTTL
		}
		c.ttl = d
		return nil
	}
}

// applyOptions applies opts to a zero-value Config and validates the result.
// Returns ErrInvalidCapacity if no WithCapacity option was provided.
func applyOptions(opts []Option) (*Config, error) {
	cfg := &Config{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if cfg.capacity == 0 {
		return nil, ErrInvalidCapacity
	}
	return cfg, nil
}
