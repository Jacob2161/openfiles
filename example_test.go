package openfiles_test

import (
	"log/slog"

	"github.com/jacob2161/openfiles"
)

func Example() {
	before, limit, err := openfiles.Raise()
	if err != nil {
		slog.Warn("raising the open file limit failed", "soft", before.Soft, "error", err)
	}
	// Each open database holds five descriptors: the file, its journal, and
	// three more.
	openDatabases := min(4096, limit.ItemsAllowed(5, openfiles.DefaultHeadroom))
	slog.Info("open file limit", "soft", limit.Soft, "open_databases", openDatabases)
}
