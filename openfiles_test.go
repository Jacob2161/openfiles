package openfiles

import (
	"errors"
	"math"
	"runtime"
	"testing"
)

// TestItemsAllowed checks the items a soft limit allows after headroom.
func TestItemsAllowed(t *testing.T) {
	for _, allowedCase := range []struct {
		limit       Limit
		perItem     int
		headroom    int
		want        int
		description string
	}{
		{Limit{Soft: 1024}, 5, 512, 102, "1,024 less 512, five each"},
		{Limit{Soft: 65536}, 4, 512, 16256, "65,536 less 512, four each"},
		{Limit{Soft: 256}, 5, 512, 1, "below the headroom still allows one"},
		{Limit{Soft: 1024}, 0, 0, 1024, "a zero count counts one"},
		{Limit{}, 5, 512, 1, "a zero soft limit still allows one"},
		{Limit{Soft: math.MaxUint64}, 5, 512, math.MaxInt, "an unlimited soft limit allows any number"},
		{Limit{Soft: 1 << 32}, 4, 0, 1 << 30, "a large limit divides before it is clamped"},
	} {
		if got := allowedCase.limit.ItemsAllowed(allowedCase.perItem, allowedCase.headroom); got != allowedCase.want {
			t.Errorf("%s: ItemsAllowed = %d, want %d", allowedCase.description, got, allowedCase.want)
		}
	}
}

// TestSpare checks descriptors left after reserved ones and headroom.
func TestSpare(t *testing.T) {
	if got := (Limit{Soft: 4096}).Spare(1000, 512); got != 2584 {
		t.Fatalf("Spare = %d, want 2584", got)
	}
	if got := (Limit{Soft: 1024}).Spare(2000, 512); got != 0 {
		t.Fatalf("Spare past the limit = %d, want 0", got)
	}
	if got := (Limit{Soft: math.MaxUint64}).Spare(0, 0); got != math.MaxInt {
		t.Fatalf("Spare of an unlimited soft limit = %d", got)
	}
	if got := (Limit{}).Spare(0, 0); got != 0 {
		t.Fatalf("Spare of a zero soft limit = %d, want 0", got)
	}
}

// TestRaise raises the real limit and checks it did not fall.
func TestRaise(t *testing.T) {
	before, after, err := Raise()
	if runtime.GOOS == "windows" {
		if !errors.Is(err, ErrUnsupported) {
			t.Fatalf("Raise on Windows: %v", err)
		}
		return
	}
	if err != nil {
		t.Logf("raising failed, as it may in a sandbox: %v", err)
	}
	if after.Soft < before.Soft || after.Hard != before.Hard {
		t.Fatalf("Raise went from %+v to %+v", before, after)
	}
	current, err := Current()
	if err != nil || current != after {
		t.Fatalf("Current = %+v, %v; Raise left %+v", current, err, after)
	}
}
