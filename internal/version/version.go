package version

import "log/slog"

var (
	Commit string
	Branch string
)

var (
	CommitAttr = slog.Attr{Key: "commit", Value: slog.StringValue(Commit)}
	BranchAttr = slog.Attr{Key: "branch", Value: slog.StringValue(Branch)}
)
