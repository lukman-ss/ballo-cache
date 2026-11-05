package inmemory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lukman-ss/ballo-cache"
)

func TestStore_SetAndGet(t *testing.T) {
	store := New()
	ctx := context.Background()

	if err := store.Set(ctx, "key", "value", 0); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := store.Get(ctx, "key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != "value" {
		t.Errorf("expected value 'value', got '%v'", val)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	store := New()
	ctx := context.Background()

	_, err := store.Get(ctx, "nonexistent")
	if err != ballo.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestStore_Delete(t *testing.T) {
	store := New()
	ctx := context.Background()

	_ = store.Set(ctx, "key", "value", 0)
	_ = store.Delete(ctx, "key")

	_, err := store.Get(ctx, "key")
	if err != ballo.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound after Delete, got %v", err)
	}
}

func TestStore_EmptyKey(t *testing.T) {
	store := New()
	ctx := context.Background()

	if _, err := store.Get(ctx, ""); err != ballo.ErrEmptyKey {
		t.Errorf("expected ErrEmptyKey on Get, got %v", err)
	}
	if err := store.Set(ctx, "", "v", 0); err != ballo.ErrEmptyKey {
		t.Errorf("expected ErrEmptyKey on Set, got %v", err)
	}
	if err := store.Delete(ctx, ""); err != ballo.ErrEmptyKey {
		t.Errorf("expected ErrEmptyKey on Delete, got %v", err)
	}
}

func TestStore_WhitespaceKey(t *testing.T) {
	store := New()
	ctx := context.Background()

	if _, err := store.Get(ctx, "   "); err != ballo.ErrEmptyKey {
		t.Errorf("expected ErrEmptyKey on whitespace key, got %v", err)
	}
}

func TestStore_Close(t *testing.T) {
	store := New()
	ctx := context.Background()

	_ = store.Set(ctx, "key", "value", 0)

	if err := store.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if _, err := store.Get(ctx, "key"); err == nil {
		t.Error("expected error after Close, got nil")
	}
	if err := store.Set(ctx, "key", "value", 0); err == nil {
		t.Error("expected error after Close, got nil")
	}
	if err := store.Delete(ctx, "key"); err == nil {
		t.Error("expected error after Close, got nil")
	}
}

func TestStore_TTLExpiration(t *testing.T) {
	store := New()
	ctx := context.Background()

	_ = store.Set(ctx, "key", "value", 1*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	_, err := store.Get(ctx, "key")
	if err != ballo.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound after TTL, got %v", err)
	}
}

func TestStore_ContextCanceled(t *testing.T) {
	store := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.Get(ctx, "key")
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}

	err = store.Set(ctx, "key", "value", 0)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled on Set, got %v", err)
	}

	err = store.Delete(ctx, "key")
	if err != context.Canceled {
		t.Errorf("expected context.Canceled on Delete, got %v", err)
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	store := New()
	ctx := context.Background()

	const routines = 100
	const opsPerRoutine = 100

	var wg sync.WaitGroup
	wg.Add(routines)

	// Concurrent writes
	for i := 0; i < routines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerRoutine; j++ {
				key := keyForID(id, j)
				_ = store.Set(ctx, key, key, 0)
			}
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	wg.Add(routines)
	for i := 0; i < routines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerRoutine; j++ {
				key := keyForID(id, j)
				val, err := store.Get(ctx, key)
				if err != nil {
					t.Errorf("Get failed for %s: %v", key, err)
					return
				}
				if val != key {
					t.Errorf("value mismatch for %s: expected %s, got %v", key, key, val)
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestStore_ConcurrentClose(t *testing.T) {
	store := New()
	ctx := context.Background()

	const routines = 50

	var wg sync.WaitGroup
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func() {
			defer wg.Done()
			_ = store.Set(ctx, "key", "value", 0)
			_, _ = store.Get(ctx, "key")
			_ = store.Delete(ctx, "key")
		}()
	}

	wg.Wait()

	if err := store.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func keyForID(id, op int) string {
	return fmt.Sprintf("k%02x%02x", id%256, op%256)
}
