package raft

import (
	"RaftLike/pkg/logger"
	"context"
	"sync"
	"time"
)

type Heartbeater struct {
	logger   *logger.Logger
	interval time.Duration
	wg       sync.WaitGroup
}

func NewHeartbeater(logger *logger.Logger, interval time.Duration) *Heartbeater {
	return &Heartbeater{logger, interval, sync.WaitGroup{}}
}

func (h Heartbeater) Run(ctx context.Context) {
	h.wg.Add(1)
	defer h.wg.Done()
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h.heartbeat()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (h Heartbeater) Stop(cancel context.CancelFunc) {
	cancel()
	h.wg.Wait()
}

func (h Heartbeater) heartbeat() {
	h.logger.Info("heartbeat")
}
