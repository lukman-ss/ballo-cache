package inmemory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lukman-ss/ballo-cache"
)

// Store is a concurrent-safe in-memory cache.
//
// Lazy cleanup: expired entries are removed by Get, Set, or Delete accesses.
// No background goroutine; some expired keys remain until accessed or the
// store is reset.
type Store struct {
	mu     sync.RWMutex
	data   map[string]entry
	closed bool
}

type entry struct {
	value     any
	expiry    time.Time
	hasExpiry bool
	gen       uint64
}

var genCounter uint64

// New returns a new Store.
func New() *Store {
	return &Store{
		data: make(map[string]entry),
	}
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
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

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	now := time.Now()

	s.mu.RLock()
	closed := s.closed
	if !closed {
		e, ok := s.data[key]
		s.mu.RUnlock()

		if !ok {
			return nil, ballo.ErrKeyNotFound
		}

		if e.hasExpiry && now.After(e.expiry) {
			s.mu.Lock()
			// Re-check under write lock: entry might have been updated
			if curr, exists := s.data[key]; exists && curr.gen == e.gen {
				delete(s.data, key)
			}
			s.mu.Unlock()
			return nil, ballo.ErrKeyNotFound
		}

		return e.value, nil
	}
	s.mu.RUnlock()

	return nil, errors.New("ballo: store closed")
}

func (s *Store) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if err := s.validateKey(key); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	now := time.Now()

	var exp time.Time
	var hasExpiry bool
	if ttl > 0 {
		exp = now.Add(ttl)
		hasExpiry = true
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return errors.New("ballo: store closed")
	}
	s.data[key] = entry{value: value, expiry: exp, hasExpiry: hasExpiry, gen: atomic.AddUint64(&genCounter, 1)}
	s.mu.Unlock()
	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	if err := s.validateKey(key); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return errors.New("ballo: store closed")
	}
	delete(s.data, key)
	s.mu.Unlock()
	return nil
}
