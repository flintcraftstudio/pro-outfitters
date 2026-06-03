// Package cleanup runs periodic reconciliation tasks in the
// background — purges expired admin sessions and consumed/expired
// magic-link tokens. Low-frequency by design (hourly default); the
// work isn't latency-sensitive and running more often just burns CPU.
// First tick fires on startup so a restart after an outage catches up
// rather than waiting a full interval.
package cleanup

import (
	"context"
	"log/slog"
	"time"

	"github.com/flintcraftstudio/standard-template/internal/store"
)

type Config struct {
	// Interval between sweeps. Defaults to 1h.
	Interval time.Duration
}

type Worker struct {
	store    *store.Store
	log      *slog.Logger
	interval time.Duration
}

func NewWorker(cfg Config, st *store.Store, log *slog.Logger) *Worker {
	if cfg.Interval <= 0 {
		cfg.Interval = 1 * time.Hour
	}
	return &Worker{
		store:    st,
		log:      log,
		interval: cfg.Interval,
	}
}

// Run blocks until ctx is done. Sweeps on startup, then every Interval.
func (w *Worker) Run(ctx context.Context) {
	w.log.Info("cleanup worker starting", "interval", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.sweep(ctx)

	for {
		select {
		case <-ctx.Done():
			w.log.Info("cleanup worker stopping")
			return
		case <-ticker.C:
			w.sweep(ctx)
		}
	}
}

func (w *Worker) sweep(ctx context.Context) {
	if err := w.store.DeleteExpiredSessions(ctx); err != nil {
		w.log.Error("cleanup: delete expired sessions", "err", err)
	}
	if err := w.store.DeleteExpiredLoginTokens(ctx); err != nil {
		w.log.Error("cleanup: delete expired login tokens", "err", err)
	}
}
