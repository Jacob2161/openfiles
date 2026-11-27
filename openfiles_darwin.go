package openfiles

import "golang.org/x/sys/unix"

// current reads RLIMIT_NOFILE.
func current() (Limit, error) {
	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		return Limit{}, err
	}
	return Limit{Soft: limit.Cur, Hard: limit.Max}, nil
}

// raise sets the soft limit to the hard one, within kern.maxfilesperproc:
// macOS reports an unlimited hard limit but refuses a soft one above that
// sysctl.
func raise() (Limit, Limit, error) {
	before, err := current()
	if err != nil {
		return Limit{}, Limit{}, err
	}
	target := before.Hard
	perProcess, err := unix.SysctlUint32("kern.maxfilesperproc")
	if err == nil && uint64(perProcess) < target {
		target = uint64(perProcess)
	}
	if before.Soft >= target {
		return before, before, nil
	}
	if err := unix.Setrlimit(unix.RLIMIT_NOFILE, &unix.Rlimit{Cur: target, Max: before.Hard}); err != nil {
		return before, before, err
	}
	after, err := current()
	if err != nil {
		return before, before, err
	}
	return before, after, nil
}
