//go:build unix && !darwin

package openfiles

import "golang.org/x/sys/unix"

// current reads RLIMIT_NOFILE.
func current() (Limit, error) {
	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		return Limit{}, err
	}
	return Limit{Soft: uint64(limit.Cur), Hard: uint64(limit.Max)}, nil
}

// raise sets the soft limit to the hard one.
func raise() (Limit, Limit, error) {
	before, err := current()
	if err != nil {
		return Limit{}, Limit{}, err
	}
	if before.Soft >= before.Hard {
		return before, before, nil
	}
	var requested unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &requested); err != nil {
		return before, before, err
	}
	requested.Cur = requested.Max
	if err := unix.Setrlimit(unix.RLIMIT_NOFILE, &requested); err != nil {
		return before, before, err
	}
	after, err := current()
	if err != nil {
		return before, before, err
	}
	return before, after, nil
}
