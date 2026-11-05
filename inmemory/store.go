package inmemory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/lukman-ss/ballo-cache"
)

// Store is a concurrent-safe in-memory cache.
//
// Lazy cleanup: expired entries are removed by Get, Set, or Delete accesses.
// No background goroutine; some expired keys remain until accessed or the
// store is reset.
type Store struct {
	mu      sync.RWMutex
	data    map[string]entry
	expires map[string]time.Time
	closed  bool
}

type entry struct {
	value any
	at    time.Time
}

// New returns a new Store.
func New() *Store {
	return &Store{
		data:    make(map[string]entry),
		expires: make(map[string]time.Time),
	}
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

func (s *Store) check() error {
	s.mu.RLock()
	closed := s.closed
	s.mu.RUnlock()
	if closed {
		return errors.New("ballo: store closed")
	}
	return nil
}

func (s *Store) validateKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return ballo.ErrEmptyKey
	}
	return nil
}

func (s *Store) Get(ctx context.Context, key string) (any, error) {
	if err := s.validateKey(key); err != nil {
		return nil, err
	}
	if err := s.check(); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	now := time.Now()

	s.mu.RLock()
	e, ok := s.data[key]
	exp, hasExpiry := s.expires[key]
	s.mu.RUnlock()

	if !ok {
		return nil, ballo.ErrKeyNotFound
	}

	if hasExpiry && now.After(exp) {
		s.mu.Lock()
		delete(s.data, key)
		delete(s.expires, key)
		s.mu.Unlock()
		return nil, ballo.ErrKeyNotFound
	}

	return e.value, nil
}

func (s *Store) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if err := s.validateKey(key); err != nil {
		return err
	}
	if err := s.check(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	now := time.Now()

	var exp time.Time
	if ttl > 0 {
		exp = now.Add(ttl)
	}

	s.mu.Lock()
	delete(s.expires, key)
	s.data[key] = entry{value: value, at: now}
	if ttl > 0 {
		s.expires[key] = exp
	}
	s.mu.Unlock()
	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	if err := s.validateKey(key); err != nil {
		return err
	}
	if err := s.check(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	delete(s.data, key)
	delete(s.expires, key)
	s.mu.Unlock()
	return nil
}
