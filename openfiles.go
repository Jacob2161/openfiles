// Package openfiles raises the process's open file limit (RLIMIT_NOFILE) and
// sizes concurrency to fit under it.
//
// The soft limit is often 1,024, which caps a server that keeps a few
// descriptors per connection, open database, or transfer long before memory
// does. Raise lifts the soft limit to the hard one, within macOS's
// kern.maxfilesperproc. ItemsAllowed and Spare then say how much the limit
// allows.
//
// On systems without RLIMIT_NOFILE, such as Windows, Current and Raise return
// ErrUnsupported.
package openfiles

import (
	"errors"
	"math"
	"runtime"
)

// DefaultHeadroom is a reasonable count of descriptors to keep free for
// everything not sized to the limit: listeners, logs, and files opened now
// and then.
const DefaultHeadroom = 512

// ErrUnsupported means the system has no RLIMIT_NOFILE.
var ErrUnsupported = errors.New("RLIMIT_NOFILE is not supported on " + runtime.GOOS)

// Limit is the process's open file limit, in descriptors.
type Limit struct {
	// Soft is the limit in force.
	Soft uint64
	// Hard is the most an unprivileged process may raise Soft to.
	Hard uint64
}

// Current returns the limit in force.
func Current() (Limit, error) {
	return current()
}

// Raise lifts the soft limit to the hard one, as far as the system allows,
// and returns the limit before and after. Go's runtime already raises it on
// Linux at start; Raise makes it explicit, and covers macOS, where the hard
// limit is often unlimited but a soft one above kern.maxfilesperproc is
// refused. When raising fails, after is before.
func Raise() (before Limit, after Limit, err error) {
	return raise()
}

// Unlimited reports whether the soft limit is the system's infinity
// (RLIM_INFINITY), which caps nothing.
func (l Limit) Unlimited() bool {
	// Systems spell infinity as the largest uint64 or the largest int64.
	return l.Soft >= math.MaxInt64
}

// ItemsAllowed returns how many items, of descriptorsPerItem descriptors each,
// the soft limit allows after headroom descriptors, at least 1, so a process
// sized by it can always do some work. An unlimited Limit allows any number:
// ItemsAllowed returns the largest int.
func (l Limit) ItemsAllowed(descriptorsPerItem int, headroom int) int {
	if l.Unlimited() {
		return math.MaxInt
	}
	spare := l.spare(0, headroom)
	items := spare / uint64(max(descriptorsPerItem, 1))
	return int(min(max(items, 1), uint64(math.MaxInt)))
}

// Spare returns the descriptors left under the soft limit after reserved and
// headroom ones, at least 0. An unlimited Limit has the largest int to spare.
func (l Limit) Spare(reserved int, headroom int) int {
	if l.Unlimited() {
		return math.MaxInt
	}
	return int(min(l.spare(reserved, headroom), uint64(math.MaxInt)))
}

// spare is Spare before it is clamped to an int.
func (l Limit) spare(reserved int, headroom int) uint64 {
	used := uint64(max(reserved, 0)) + uint64(max(headroom, 0))
	if l.Soft <= used {
		return 0
	}
	return l.Soft - used
}
