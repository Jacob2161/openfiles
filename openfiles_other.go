//go:build !unix

package openfiles

func current() (Limit, error) { return Limit{}, ErrUnsupported }

func raise() (Limit, Limit, error) { return Limit{}, Limit{}, ErrUnsupported }
