package ballo

import (
	"context"
	"errors"
	"time"
)

// Store is a minimal key-value cache with optional TTL and context support.
//
// Safety: implementations MUST be safe for concurrent use.
//
// Key expired behavior: a key whose TTL has elapsed is treated as missing
// (Get returns ErrKeyNotFound). Expired keys may linger in memory until the
// next access or until an optional background sweep removes them. No
// guarantee is made about exact expiration timing; TTL is a best-effort upper
// bound on how long a value remains readable.
type Store interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

// ErrKeyNotFound is returned when the requested key does not exist or has expired.
var ErrKeyNotFound = errors.New("ballo: key not found")

// ErrEmptyKey is returned when a key is empty.
var ErrEmptyKey = errors.New("ballo: key is empty") // ponytail: keep for backward compatibility
