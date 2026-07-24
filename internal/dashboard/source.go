package dashboard

import (
	"context"
	"time"

	"github.com/setup-env/app/internal/system"
)

const (
	DefaultRefreshInterval    = time.Second
	DefaultFilesystemInterval = 5 * time.Second
	DefaultDiagnosticInterval = time.Minute
	DefaultHistoryLimit       = 60
)

type RefreshRequest struct {
	Previous           system.Snapshot
	CollectedAt        time.Time
	IncludeFilesystems bool
	IncludeDiagnostics bool
}

type Source interface {
	Initial(context.Context) (system.Snapshot, error)
	Refresh(context.Context, RefreshRequest) (system.Snapshot, error)
}
