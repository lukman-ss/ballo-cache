# ballo-cache

## Installation

```bash
 go get github.com/lukman-ss/ballo-cache
```

## Usage example

```go
package main

import (
 "context"
 "fmt"
 "time"

 "github.com/lukman-ss/ballo-cache"
 "github.com/lukman-ss/ballo-cache/inmemory"
)

func main() {
 ctx := context.Background()
 store := inmemory.New()
 _ = store.Set(ctx, "hello", "world", 5*time.Second)
 val, err := store.Get(ctx, "hello")
 fmt.Println(val, err)
 _ = store.Close()
}
```

## API

`ballo.Store` is the minimal cache interface: `Get`, `Set`, `Delete`, and
`Close`. Each method takes a `context.Context`; context cancellation fails the
operation immediately.

- Empty key returns `ballo.ErrEmptyKey`.
- Missing or expired key returns `ballo.ErrKeyNotFound` from `Get`.
- Expired keys are removed lazily on access; no background sweep runs.
- `Close` marks the store closed; subsequent operations return an error until
  a new store is created.

## Prerequisites

- Go 1.23+ (module aware)
- No external services; in‑memory only.

## Limitations

- No persistence; data disappears when process exits.
- No background eviction; expired entries are removed lazily on access.
- Only `context.Context` cancellation is respected; deadlines are not used for TTL.

### Thread-safety notes

- `Get` reads value and expiry atomically; prevents reading stale value then deleting fresh entry
- All operations check closed status under the same lock as `Close`; no operations run after `Close` returns
- TTL expiry re-checks entry identity under write lock before deletion
- Entry identity uses atomic generation counter (avoids comparing slice/map values)

## Running tests

```bash
 go test ./... -race
```

## Change Log

### Unreleased

### Added
- `Store` interface with `Get`, `Set`, `Delete`, and `Close` and context support.
- `ErrKeyNotFound` for missing or expired keys; `ErrEmptyKey` for empty keys.
- `inmemory.Store` implementation (concurrent-safe, lazy expiration, `Close`).
- Usage documentation and tests including TTL, empty key, closed store, and race tests.

