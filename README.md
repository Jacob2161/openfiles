# openfiles

[![Go Reference](https://pkg.go.dev/badge/github.com/jacob2161/openfiles.svg)](https://pkg.go.dev/github.com/jacob2161/openfiles)

Raise the process's open file limit (RLIMIT_NOFILE) to what the system allows, including macOS's
per-process cap, and size concurrency to fit under it.

```bash
go get github.com/jacob2161/openfiles
```

## Usage

```go
before, limit, err := openfiles.Raise()
if err != nil {
	log.Printf("raising the open file limit from %d failed: %v", before.Soft, err)
}
// Each open database holds five descriptors.
openDatabases := min(configured, limit.ItemsAllowed(5, openfiles.DefaultHeadroom))
connections := limit.Spare(openDatabases*5, openfiles.DefaultHeadroom)
```

- `Raise` lifts the soft limit to the hard one. On macOS, where the hard limit
  is often unlimited but a soft limit above `kern.maxfilesperproc` is refused,
  it stops at that sysctl.
- `ItemsAllowed` and `Spare` size concurrency under the limit in force; an unlimited soft
  limit (RLIM_INFINITY) allows any number.
- On Windows, `Current` and `Raise` return `ErrUnsupported`.

Requires Go 1.27. Licensed under the MIT license.
