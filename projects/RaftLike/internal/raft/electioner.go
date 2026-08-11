package raft

import (
	"RaftLike/pkg/logger"
	"context"
	"sync"
	"time"
)

type Electioner struct {
	logger   *logger.Logger
	interval time.Duration
	wg       sync.WaitGroup
}

func NewElectioner(logger *logger.Logger, interval time.Duration) *Electioner {
	return &Electioner{logger, interval, sync.WaitGroup{}}
}

func (e Electioner) Run(ctx context.Context) {
	e.wg.Add(1)
	defer e.wg.Done()
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			e.elect()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (h Electioner) Stop(cancel context.CancelFunc) {
	cancel()
	h.wg.Wait()
}

func (h Electioner) elect() {
	h.logger.Info("elect")
}
