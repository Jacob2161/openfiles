# Changelog

This module follows semantic versioning. Before v1, a minor version may
change the API; each such change says what callers must do.

## Unreleased

- `Limit.Soft` and `Limit.Hard` each have a doc comment. No API change.
- CI also checks formatting, `go mod tidy -diff`, and `govulncheck`.

## v0.1.0 (2026-11-27)

First release.

- `Raise` lifts the open file soft limit to the hard one, within macOS's
  `kern.maxfilesperproc`, and `Current` reads it.
- `Limit.ItemsAllowed` and `Limit.Spare` size concurrency to fit under the
  limit after headroom.
- On systems without RLIMIT_NOFILE, such as Windows, `Current` and `Raise`
  return `ErrUnsupported`.
