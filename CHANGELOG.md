# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Fix concurrent race in `Get`: read expiry and value atomically from single `entry` struct to prevent reading stale value then deleting fresh entry
- Fix `check()` closed status race: guard all operations with same lock as `Close`
- Add `entry.expiry`, `entry.hasExpiry` to store TTL metadata inline
- Update TTL expiry cleanup to re-check entry identity under write lock

### Added

- `TestStore_ConcurrentGetSetSameKey`: verify Get/Set on shared key under contention
- `TestStore_OperationVsClose`: verify operations fail after Close
- `TestStore_ConcurrentGetSetDeleteClose`: stress-test mixed ops with concurrent Close

