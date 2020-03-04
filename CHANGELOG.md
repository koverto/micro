# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog][], and this project adheres to
[Semantic Versioning][].

## Unreleased

## v1.1.0 - 2020-03-04

### Added

- `ContextWithRequestID(ctx context.Context, rid *uuid.UUID) context.Context`
- `RequestIDFromContext(ctx context.Context) *uuid.UUID`

### Changed

- Updated github.com/micro/go-micro/v2 from 2.1.2 to 2.2.0

## v1.0.1 - 2020-02-25

### Changed

- Updated gihtub.com/micro/go-micro/v2 from 2.1.1 to 2.1.2

## v1.0.0 - 2020-02-21

- Initial release

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html
