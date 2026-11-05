# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `Store` interface with `Get`, `Set`, `Delete`, and `Close` and context support.
- `ErrKeyNotFound` for missing or expired keys; `ErrEmptyKey` for empty keys.
- `inmemory.Store` implementation (concurrent-safe, lazy expiration, `Close`).
- Tests covering basic ops, TTL expiration, empty key, closed store, and concurrent access under `-race`.

