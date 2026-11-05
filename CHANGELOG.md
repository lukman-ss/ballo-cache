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
- Fixed expiration cleanup panic for slice/map values by using generation counter instead of value comparison
- Updated tests for slice/map expiration and concurrent updates
- Refactored entry struct with generation ID
- Added atomic counter for entry generations
- Adjusted entry deletion check to compare generation only

### Added

- `TestStore_ExpiredSliceValue`: verify expiration works with slice values
- `TestStore_ExpiredMapValue`: verify expiration works with map values
- `TestStore_ExpiredSliceValueConcurrentUpdate`: verify Get cleanup doesn't delete updated slice entries
- `TestStore_ExpiredMapValueConcurrentUpdate`: verify Get cleanup doesn't delete updated map entries
